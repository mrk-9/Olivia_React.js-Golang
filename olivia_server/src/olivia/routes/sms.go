package routes

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"olivia/models"
	//"olivia/twilio"
	"strconv"
	"olivia/twilio"
	"olivia/zendesk"
)

func ReceiveSMS(ctx *fasthttp.RequestCtx) {
	type SMSRequest struct {
		Sender   string `json:"From"`
		Receiver string `json:"To"`
		Message  string `json:"Body"`
	}

	type WSMessage struct {
		Type string `json:"type"`
		Message models.Message `json:"message"`
	}

	var smsRequest SMSRequest

	smsRequest.Sender = string(ctx.PostArgs().Peek("From")[:])
	smsRequest.Receiver = string(ctx.PostArgs().Peek("To")[:])
	smsRequest.Message = string(ctx.PostArgs().Peek("Body")[:])

	var lead models.Lead
	lead.Phone = smsRequest.Sender
	lead.FindByPhone()

	if lead.Id == 0 {
		log.Printf("Unknown sender %s", smsRequest.Sender)
		return
	}

	if lead.ZendeskTicketId == 0 {
		// No Zendesk ticket - create one
		ticket := zendesk.CreateTicket(lead.Name, fmt.Sprintf("Lead ID: %d\r\nPhone number: %s\r\nMessage: %s", lead.Id, lead.Phone, smsRequest.Message))
		lead.SetZendeskTicketId(ticket.Id)
	} else {
		var ticket zendesk.Ticket
		ticket.Id = lead.ZendeskTicketId
		ticket.Comment.Body = fmt.Sprintf("Message: %s", smsRequest.Message)
		ticket.Update()
	}

	var user models.User
	user.GetByInternalPhoneNumber(smsRequest.Receiver)

	if user.Id == 0 {
		log.Printf("Unknown receiver %s", smsRequest.Receiver)
		return
	}

	var message models.Message
	message.ForLead = false
	message.Message = smsRequest.Message
	message.Lead = lead.Id
	message.User = user.Id
	message.Create()

	var wsMessage WSMessage
	wsMessage.Type = "message"
	wsMessage.Message = message
	SendWSToUser(user, wsMessage)

	fmt.Fprintf(ctx, "")
}

func SendSMS(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/json")

	type SendSMSRequest struct {
		Lead    int64    `json:"lead"`
		Message string `json:"message"`
	}

	var sendSMSRequest SendSMSRequest
	err := json.Unmarshal(ctx.Request.Body(), &sendSMSRequest)
	if err != nil {
		log.Printf("Error parsing SendSMS request: %s", err.Error())
	}

	userInterface := ctx.UserValue("user")
	user := userInterface.(models.User)

	pn, pnErr := user.GetPhoneNumber()
	if pnErr != nil {
		log.Println(pnErr.Error())
	}

	if pn.Id == 0 {
		log.Printf("Send number not found. %v", pn)
	}

	var lead models.Lead
	lead.GetById(sendSMSRequest.Lead)

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

func GetChat(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/json")

	leadId, _ := strconv.ParseInt(string(ctx.QueryArgs().Peek("lead")[:]), 10, 64)

	userInterface := ctx.UserValue("user")
	user := userInterface.(models.User)

	var lead models.Lead
	lead.GetById(leadId)

	messages, _ := models.GetChat(user, lead)

	RespondJSON(200, ctx, messages)
}
