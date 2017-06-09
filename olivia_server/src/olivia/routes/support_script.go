package routes

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
	"olivia/models"
	"strconv"
)

func GetScriptsAndAnswers(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/json")

	userId, _ := strconv.ParseInt(string(ctx.QueryArgs().Peek("user")[:]), 10, 64)
	leadId, _ := strconv.ParseInt(string(ctx.QueryArgs().Peek("lead")[:]), 10, 64)

	//userInterface := ctx.UserValue("user")
	//user := userInterface.(models.SupportUser)

	var lead models.Lead
	lead.GetById(leadId)

	var realtor models.User
	realtor.Id = userId
	realtor.GetById()

	if lead.Id == 0 {
		RespondJSON(fasthttp.StatusOK, ctx, Response{"Lead ID is invalid"})
		return
	}

	if realtor.Id == 0 {
		RespondJSON(fasthttp.StatusOK, ctx, Response{"User ID is invalid"})
		return
	}

	_, property := lead.GetMostRecentProperty()

	scripts := models.GetScriptsAndAnswers(realtor, lead, property)

	RespondJSON(fasthttp.StatusOK, ctx, scripts)
}

func AnswerScript(ctx *fasthttp.RequestCtx) {
	type AnswerScriptRequest struct {
		ScriptId int64  `json:"scriptId"`
		LeadId   int64  `json:"lead"`
		Text     string `json:"answer"`
	}

	var request AnswerScriptRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
		log.Println(err.Error())
	}

	_, scriptAnswer := models.GetAnswerByScriptAndLead(request.ScriptId, request.LeadId)

	if scriptAnswer.Id == 0 {
		scriptAnswer.Lead = request.LeadId
		scriptAnswer.Script = request.ScriptId
		scriptAnswer.Text = request.Text
		scriptAnswer.Create()
	} else {
		scriptAnswer.Text = request.Text
		scriptAnswer.Save()
	}

	ctx.SetContentType("text/json")

	RespondJSON(fasthttp.StatusOK, ctx, Response{})
}
