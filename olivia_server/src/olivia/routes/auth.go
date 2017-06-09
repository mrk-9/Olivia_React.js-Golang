package routes

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
	"olivia/models"
	"olivia/twilio"
	"strings"
	"fmt"
)

func Register(ctx *fasthttp.RequestCtx) {
	type RegisterRequest struct {
		FirstName  string `json:"fname"`
		SecondName string `json:"sname"`
		Email      string `json:"email"`
		Code       string `json:"code"`
		Phone      string `json:"number"`
		Password   string `json:"password"`
	}

	var registerRequest RegisterRequest
	err := json.Unmarshal(ctx.Request.Body(), &registerRequest)
	if err != nil {
	}

	ctx.SetContentType("text/json")

	defer ctx.Request.ConnectionClose()

	if registerRequest.FirstName == "" {
		RespondJSON(400, ctx, Response{"First name should be present"})
		return
	}

	if registerRequest.SecondName == "" {
		RespondJSON(400, ctx, Response{"Second name should be present"})
		return
	}

	if registerRequest.Email == "" {
		RespondJSON(400, ctx, Response{"Email should be present"})
		return
	}

	if registerRequest.Phone == "" || registerRequest.Code == "" {
		RespondJSON(400, ctx, Response{"Phone should be present"})
		return
	}

	if registerRequest.Password == "" {
		RespondJSON(400, ctx, Response{"Password should be present"})
		return
	}

	userExist, _ := models.User{}.CheckExist(registerRequest.Email)

	if userExist == true {
		RespondJSON(400, ctx, Response{"User with this email already exist"})
		return
	}

	var gluedPhoneNumber = strings.Join([]string{registerRequest.Code, registerRequest.Phone}, "")
	lookupResponse := twilio.Lookup(gluedPhoneNumber)

	var code string = registerRequest.Code
	var phone string

	phone = lookupResponse.PhoneNumber

	var user models.User
	user.Email = registerRequest.Email
	user.Code = code
	user.Phone = phone
	user.Password = registerRequest.Password
	user.FirstName = registerRequest.FirstName
	user.SecondName = registerRequest.SecondName

	savingErr := user.Save()

	if savingErr == nil {
		// Assign phone number to user
		// Lookup is falling back to US in case of failure
		var numbers = twilio.GetPhoneNumbers(lookupResponse.CountryCode)

		if len(numbers) > 0 {
			number := twilio.ConnectPhoneNumber(user, numbers[0])

			var pn models.PhoneNumber
			pn.Number = number.PhoneNumber
			pn.Sid = number.SID
			pn.Owner = user
			pn.Add()
		} else {
			log.Println(lookupResponse)
			log.Println(numbers)
		}

		// Send initial messages from lead #20
		var m1 models.Message
		m1.User = user.Id
		m1.Lead = 20
		m1.ForLead = false
		m1.Message = fmt.Sprintf("Hi %s, welcome to your free 10 day trial! Olivia is a human assistant that engages your incoming leads 24/7. We communicate with your leads via SMS within minutes of their enquiry, ensuring a response rate that exceeds email or phone communications.", user.FirstName)
		m1.Create()

		//var m2 models.Message
		//m2.User = user.Id
		//m2.Lead = 20
		//m2.ForLead = false
		//m2.Message = "Olivia agents qualify your prospects using a script you approve, letting you focus on selling. You can <a href=\"/scripts\" target=\"_blank\">create your own script</a> for our agents to use, or start with our default script."
		//m2.Create()

		var m3 models.Message
		m3.User = user.Id
		m3.Lead = 20
		m3.ForLead = false
		m3.Message = "When a prospect sends an enquiry, Olivia agents qualify them using a script you approve, which lets you focus on selling. You can <a href=\"/scripts\">create your own script</a> for our agents to use, or start with our default script."
		m3.Create()

		var m4 models.Message
		m4.User = user.Id
		m4.Lead = 20
		m4.ForLead = false
		m4.Message = fmt.Sprintf("Your private Olivia email address is <strong>%s@%s</strong>. In order for Olivia’s agents to qualify your prospects, you’ll need to set up your email inbox to automatically forward emails from listing sites to your private Olivia email address.", user.OEmail, "oliviaemail.com")
		m4.Create()

		var m5 models.Message
		m5.User = user.Id
		m5.Lead = 20
		m5.ForLead = false
		m5.Message = "If you need help, you can read our email set up instructions, or tap the chat box on the lower right for assistance."
		m5.Create()

		models.CreateDefaultScripts(user.Id)

		RespondJSON(200, ctx, Response{"Successfully registered"})
	} else {
		RespondJSON(400, ctx, savingErr)
	}
}

func Login(ctx *fasthttp.RequestCtx) {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequest
	err := json.Unmarshal(ctx.Request.Body(), &loginRequest)
	if err != nil {
	}

	ctx.SetContentType("text/json")

	defer ctx.Request.ConnectionClose()

	if loginRequest.Email == "" {
		RespondJSON(400, ctx, Response{"Email should be present"})
		return
	}

	if loginRequest.Password == "" {
		RespondJSON(400, ctx, Response{"Password should be present"})
		return
	}

	var user models.User
	user.GetByCredentials(loginRequest.Email, loginRequest.Password)

	if user.Id == 0 {
		RespondJSON(400, ctx, Response{"User not found"})
		return
	}

	var token models.Token
	token.GetTokenForUser(user)

	if token.Id == 0 {
		token.Create(user)
	}

	RespondJSON(200, ctx, token)
}

func GetProfile(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/json")

	type AssignedPhoneNumber struct {
		AssignedPhoneNumber string `json:"assigned_phone_number"`
	}

	type Response struct {
		models.User
		AssignedPhoneNumber
	}

	userInterface := ctx.UserValue("user")
	user := userInterface.(models.User)

	twilioPhoneNumber, _ := user.GetPhoneNumber()

	assignedPhoneNumber := AssignedPhoneNumber{twilioPhoneNumber.Number}

	var response = Response{user, assignedPhoneNumber}

	RespondJSON(200, ctx, response)
}

func ValidatePhone(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/json")

	type LookupRequest struct {
		Code   string `json:"code"`
		Number string `json:"number"`
	}

	var request LookupRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		log.Println(err.Error())
	}

	result := twilio.Lookup(strings.Join([]string{request.Code, request.Number}, ""))

	isValid := false
	if result.PhoneNumber != "" {
		isValid = true
	}

	RespondJSON(fasthttp.StatusOK, ctx, map[string]bool{"valid": isValid})
}
