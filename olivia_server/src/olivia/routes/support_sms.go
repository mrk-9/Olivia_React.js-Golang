package routes

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
	"olivia/models"
	"olivia/twilio"
	"strconv"
)

func GetSupportChat(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/json")

	userId, _ := strconv.ParseInt(string(ctx.QueryArgs().Peek("user")[:]), 10, 64)
	leadId, _ := strconv.ParseInt(string(ctx.QueryArgs().Peek("lead")[:]), 10, 64)

	//userInterface := ctx.UserValue("user")
	//user := userInterface.(models.SupportUser)

	var lead models.Lead
	lead.GetById(leadId)

	var user models.User
	user.Id = userId
	user.GetById()

	messages, _ := models.GetChat(user, lead)

	RespondJSON(200, ctx, messages)
}

func SendSupportSMS(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/json")

	type SendSMSRequest struct {
		LeadId    int64  `json:"lead"`
		RealtorId int64  `json:"user"`
		Message   string `json:"message"`
	}

	var sendSMSRequest SendSMSRequest
	err := json.Unmarshal(ctx.Request.Body(), &sendSMSRequest)
	if err != nil {
		log.Printf("Error parsing SendSMS request: %s", err.Error())
	}

	var user models.User
	user.Id = sendSMSRequest.RealtorId
	user.GetById()

	pn, pnErr := user.GetPhoneNumber()
	if pnErr != nil {
		log.Println(pnErr.Error())
	}

	if pn.Id == 0 {
		log.Printf("Send number not found. %v", pn)
	}

	var lead models.Lead
	lead.GetById(sendSMSRequest.LeadId)

	// Create message record
	var message models.Message
	message.Lead = lead.Id
	message.User = user.Id
	message.Message = sendSMSRequest.Message
	message.ForLead = true
	message.Create()

	// Send SMS
	twilio.SendMessage(lead.Phone, pn.Number, sendSMSRequest.Message)

	RespondJSON(200, ctx, Response{"Sent"})
}
