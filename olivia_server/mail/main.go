package main

import (
	"github.com/mhale/smtpd"
	"log"
	"net"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	emailParser "github.com/veqryn/go-email/email"
	"olivia/models"
	"olivia/parsers"
	"strings"
	"olivia/twilio"
	"regexp"
	"time"
	"config"
)

var DB *sql.DB

func HandleErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func parseGmailForwardingConfirmation(user *models.User, email *models.Email) {
	codeReg := regexp.MustCompile(`#([0-9]+)`)
	linkReg := regexp.MustCompile(`(https://.*)\s`)

	codeResult := codeReg.FindAllStringSubmatch(email.Subject, -1)
	if len(codeResult) == 0 {
		log.Println("Error parsing Gmail Forwarding Confirmation - 1")
		log.Println(email)
		return
	}
	code := codeResult[0][1]

	linkResult := linkReg.FindAllString(email.Text, -1)
	if len(linkResult) == 0 {
		log.Println("Error parsing Gmail Forwarding Confirmation - 2")
		log.Println(email.Text)
		return
	}
	link := linkResult[0]

	var m1 models.Message
	m1.Lead = 20
	m1.User = user.Id
	m1.ForLead = false
	m1.Message = fmt.Sprintf("To allow Google to forward emails to Olivia so that we can respond to your prospects, please click <a href=\"%s\" target=\"_blank\">this link</a>, or enter this code into your Gmail settings: %s", link, code)
	m1.Create()

	var m2 models.Message
	m2.Lead = 20
	m2.User = user.Id
	m2.ForLead = false
	m2.Message = "Once you’ve done this, you can create a filter to forward your lead emails to Olivia."
	m2.Create()
}

func mailHandler(origin net.Addr, from string, to []string, data []byte) {
	log.Printf("Received mail from %s for %s", from, to[0])

	// check receivers list to make sure we have corresponding user
	for _, emailAddress := range to {
		o_email := strings.Split(emailAddress, "@")[0]

		var user models.User
		user.OEmail = o_email
		user.GetByOEmail()

		// if we have user with this email
		if user.Id != 0 {
			log.Printf("Email for %s", user.FirstName)

			var email models.Email

			email.UserId = user.Id
			email.Raw = string(data[:])

			emailReader := strings.NewReader(email.Raw)
			parsedEmail, parseErr := emailParser.ParseMessage(emailReader)

			if parseErr != nil {
				log.Printf("Unable to parse message")
				log.Printf(parseErr.Error())
				return
			}

			email.From = from

			for _, part := range parsedEmail.MessagesAll() {
				mediaType, _, _ := part.Header.ContentType()

				switch mediaType {
				case "text/html":
					email.HTML = fmt.Sprintf("%s", part.Body)
				case "text/plain":
					email.Text = fmt.Sprintf("%s", part.Body)
				}
			}

			email.Subject = parsedEmail.Header.Subject()

			// Test whether this email is Gmail forwarding confirmation
			reg := regexp.MustCompile(`\(#[0-9]+\) Gmail Forwarding Confirmation`)
			gmailConfResult := reg.FindAllString(email.Subject, -1)
			if len(gmailConfResult) > 0 {
				parseGmailForwardingConfirmation(&user, &email)

				return
			}

			var emailData parsers.EmailData
			parsers.Parse(from, email.HTML, &emailData)

			if emailData.IsFailed() == true {
				var be models.BadEmail
				be.Text = email.Text
				be.HTML = email.Text
				be.Subject = email.Subject
				be.From = email.From
				be.Create()
				return
			}

			var leadNeedsAssistance bool = false;

			// Validate lead phone number with Twilio
			if (emailData.Phone != "") {
				log.Println("Lead phone number exist")
				lookupResult := twilio.Lookup(emailData.Phone)

				if (lookupResult.PhoneNumber != "") {
					log.Printf("Lookup work fine, setting %s", lookupResult.PhoneNumber)
					emailData.Phone = lookupResult.PhoneNumber
				} else {
					log.Println("Lookup failed, adding code to original number: %s", strings.Join([]string{user.Code, emailData.Phone}, ""))
					// If phone number is not valid - merge it with customer's code
					lookupResult = twilio.Lookup(strings.Join([]string{user.Code, emailData.Phone}, ""))

					if (lookupResult.PhoneNumber != "") {
						log.Println("Lookup succeeded: %s", lookupResult.PhoneNumber)
						emailData.Phone = lookupResult.PhoneNumber
					} else {
						log.Println("Lookup failed, needs assistance")
						// If code + phone number can't be validated too - set special flag for lead
						leadNeedsAssistance = true
					}
				}
			} else {
				leadNeedsAssistance = true
			}

			email.Status = emailData.Status
			email.EnquiryType = emailData.EnquiryType

			_, err := email.Create()
			if err != nil {
				log.Printf(err.Error())
			}

			var lead models.Lead
			lead.Name = emailData.Name
			lead.Phone = emailData.Phone
			lead.Email = emailData.EmailAddress
			lead.NeedsAssistance = leadNeedsAssistance

			// Trying to find existing lead
			lead.Search()

			// If lead not found - create him
			if lead.Id == 0 {
				log.Println("Lead not found, creating new one...")
				_, createLeadErr := lead.Create()

				if createLeadErr != nil {
					log.Println(createLeadErr.Error())
				}
			}

			email.Lead = int(lead.Id)
			email.Save()

			log.Println("Email data")
			log.Println(emailData)

			var property models.Property
			property.Address = emailData.Address
			property.Description = emailData.Description
			property.ImageUrl = emailData.ImageURL
			property.Link = emailData.Link
			property.Price = emailData.Price

			// Trying to find existing property
			property.Search()

			// If property not found - create it
			if property.Id == 0 {
				log.Println("Property not found, creating new one...")
				property.Create()
			}

			go SendInitialScript(&user, &lead, &property, 3 * 60 * 1000)

			lead.AddInterest(property)
		}
	}
}

func SendInitialScript(user *models.User, lead *models.Lead, property *models.Property, delay int) {
	time.Sleep(time.Duration(delay) * time.Millisecond)

	var m models.Message
	m.Lead = lead.Id
	m.User = user.Id
	m.ForLead = true

	pn, pnErr := user.GetPhoneNumber()
	if pnErr != nil {
		log.Println(pnErr.Error())
	}

	scripts, _ := models.GetUserScripts(*user)
	if len(scripts) == 0 {
		log.Printf("No scripts for %s", user.Id)
		return
	}

	text := models.FormatScriptMessage(scripts[0].Text, lead, property)

	m.Message = text;
	m.Create()

	log.Printf("Sending SMS to lead %s: %s", lead.Id, text)

	// Send SMS
	twilio.SendMessage(lead.Phone, pn.Number, text)
}

func main() {
	log.Printf("Starting...")

	DB, connErr := sql.Open("mysql", config.MysqlHost)
	HandleErr(connErr)

	DB.Exec("SET NAMES `utf8`")
	DB.Exec("SET CHARACTER SET utf8")

	models.DB = DB

	log.Printf("Connected to DB")
	err := smtpd.ListenAndServe("0.0.0.0:25", mailHandler, "Olivia", "")

	if err != nil {
		log.Fatal(err)
	}
}
