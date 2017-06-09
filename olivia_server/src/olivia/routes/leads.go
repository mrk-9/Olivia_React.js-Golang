package routes

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
	"olivia/models"
)

func GetLeads(ctx *fasthttp.RequestCtx) {
	type ResponseItem struct {
		models.Lead    `json:"lead"`
		models.Message `json:"last_message"`
		models.Email   `json:"email"`
	}

	ctx.SetContentType("text/json")

	userInterface := ctx.UserValue("user")
	user := userInterface.(models.User)

	var leads []models.Lead
	leads, _ = models.GetLeadsByRealtorOrFree(user.Id)

	for index := range leads {
		leads[index].GetProperties()
	}

	var response = make([]ResponseItem, len(leads))

	for index, lead := range leads {
		lastMessage, lmErr := models.GetLastMessageFromChat(user, lead)
		if lmErr != nil {
			//log.Println(lmErr.Error())
		}

		email, emailErr := models.GetEmailByLead(lead)
		if emailErr != nil {
			//log.Println(emailErr.Error())
		}

		var responseItem = ResponseItem{lead, lastMessage, email}

		response[index] = responseItem
	}

	RespondJSON(200, ctx, response)
}

func GetSupportLeads(ctx *fasthttp.RequestCtx) {
	type ResponseItem struct {
		models.Lead    `json:"lead"`
		models.Message `json:"last_message"`
		models.Email   `json:"email"`
		models.User    `json:"realtor"`
	}

	ctx.SetContentType("text/json")

	//userInterface := ctx.UserValue("user")
	//user := userInterface.(models.SupportUser)

	var leads []models.Lead
	leads, _ = models.GetFreeLeads()

	for index := range leads {
		leads[index].GetProperties()
	}

	var response = make([]ResponseItem, len(leads))

	for index, lead := range leads {
		user := lead.GetRealtorFromEmail()

		lastMessage, lmErr := models.GetLastMessageFromChat(user, lead)
		if lmErr != nil {
			//log.Println(lmErr.Error())
		}

		email, emailErr := models.GetEmailByLead(lead)
		if emailErr != nil {
			//log.Println(emailErr.Error())
		}

		var responseItem = ResponseItem{lead, lastMessage, email, user}

		response[index] = responseItem
	}

	RespondJSON(200, ctx, response)
}

func SetOwnership(ctx *fasthttp.RequestCtx) {
	type Request struct {
		Lead int64 `json:"lead"`
		Flag bool  `json:"flag"`
	}

	var request Request
	err := json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		log.Printf("Error parsing request: %s", err.Error())
	}

	userInterface := ctx.UserValue("user")
	user := userInterface.(models.User)

	type Response struct {
		Owning bool `json:"owning"`
	}

	var lead models.Lead
	lead.GetById(request.Lead)

	isYours := lead.EnsureLeadCameToUser(&user)

	if isYours != true {
		log.Println("This lead does not belong to you")
		response := Response{false}
		RespondJSON(fasthttp.StatusOK, ctx, response)
		return
	}

	if request.Flag {
		log.Println("set")
		lead.SetOwnership(&user)
	} else {
		log.Println("release")
		lead.ReleaseOwnership()
	}

	refreshErr := lead.Refresh()
	if refreshErr != nil {
		log.Println(refreshErr.Error())
	}
	log.Println(lead.Realtor)

	response := Response{Owning: (lead.Realtor == user.Id)}
	RespondJSON(fasthttp.StatusOK, ctx, response)
}

func SetAssistance(ctx *fasthttp.RequestCtx) {
	type Request struct {
		Lead int64 `json:"lead"`
		Flag bool  `json:"flag"`
	}

	var request Request
	err := json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		log.Printf("Error parsing request: %s", err.Error())
	}

	type Response struct {
		NeedsAssistance bool `json:"assistance"`
	}

	var lead models.Lead
	lead.GetById(request.Lead)

	if request.Flag {
		log.Println("assitance")
		lead.SetAssistance(true)
	} else {
		log.Println("no assistance")
		lead.SetAssistance(false)
	}

	refreshErr := lead.Refresh()
	if refreshErr != nil {
		log.Println(refreshErr.Error())
	}
	log.Println(lead.NeedsAssistance)

	response := Response{NeedsAssistance: (lead.NeedsAssistance)}
	RespondJSON(fasthttp.StatusOK, ctx, response)
}
