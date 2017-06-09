package main

import "database/sql"
import (
	_ "github.com/go-sql-driver/mysql"
	"config"
	"olivia/models"
	"log"
	//"gopkg.in/mailgun/mailgun-go.v1"
	//"fmt"
	"html/template"
	"io/ioutil"
	"bytes"
	"gopkg.in/mailgun/mailgun-go.v1"
	"fmt"
	"strings"
)

var DB *sql.DB

func HandleErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	var err error
	DB, err = sql.Open("mysql", config.MysqlHost)
	HandleErr(err)

	models.DB = DB

	DB.Exec("SET NAMES `utf8`")
	DB.Exec("SET CHARACTER SET utf8")

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	html, readFileErr := ioutil.ReadFile("./cronjob/emails/dailyrealtoremail_converted.html")
	if readFileErr != nil {
		log.Println(readFileErr)
	}

	tpl, tplErr := template.New("dailyrealtoremail").Parse(string(html))
	if tplErr != nil {
		log.Println(tplErr)
		return
	}

	type User struct {
		ID int64
		Initials string
		Name string
		Property string
	}

	type Data struct {
		NeedsAssistance []User
		SlippingAway []User
		NewProspects []User
	}

	users := models.GetAllUsers()

	convertLeads := func (leads []models.Lead) []User {
		var users = make([]User, 0)

		for _, lead := range leads {
			var user User

			user.ID = lead.Id
			user.Name = lead.Name

			splitName := strings.Split(user.Name, " ")
			user.Initials = strings.ToUpper(string(user.Name[0]))
			if len(splitName) >= 2 {
				user.Initials = strings.ToUpper(string(splitName[0][0])) + strings.ToUpper(string(splitName[1][0]))
			}

			var _, latestProperty = lead.GetMostRecentProperty()
			user.Property = latestProperty.Address

			users = append(users, user)
		}

		return users
	}

	for _, user := range users {
		var data Data

		if user.Email != "lavavrik@yandex.ru" && user.Email != "chris@getolivia.co" && user.Email != "chris@knappick.com" {
			continue
		}

		data.NeedsAssistance = convertLeads(models.GetLeadsWhoNeedAssistanceByRealtor(user.Id))
		data.SlippingAway = convertLeads(models.GetSlippingAwayLeadsByRealtor(user.Id))
		data.NewProspects = convertLeads(models.GetNewLeadsByRealtor(user.Id))

		if len(data.NeedsAssistance) == 0 && len(data.SlippingAway) == 0 && len(data.NewProspects) == 0 {
			continue
		}

		log.Println(user.Id)

		var buf = new(bytes.Buffer)

		execErr := tpl.Execute(buf, data)
		if execErr != nil {
			log.Println(execErr)
		}

		mg := mailgun.NewMailgun("getolivia.co", "key-392aceadb7b57481d125944b97c487e5", "key-392aceadb7b57481d125944b97c487e5")
		message := mailgun.NewMessage(
			"Olivia <support@getolivia.co>",
			"Daily Realtor's Email",
			"",
			user.Email)
		message.SetHtml(buf.String())

		resp, id, err := mg.Send(message)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %s Resp: %s\n", id, resp)
	}
}
