package routes

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
	"olivia/models"
	"strconv"
)

func CreateScript(ctx *fasthttp.RequestCtx) {
	type CreateScriptRequest struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	}

	var createScriptRequest CreateScriptRequest
	err := json.Unmarshal(ctx.Request.Body(), &createScriptRequest)
	if err != nil {
	}

	ctx.SetContentType("text/json")

	userInterface := ctx.UserValue("user")
	user := userInterface.(models.User)

	defer ctx.Request.ConnectionClose()

	var script models.Script
	script.Title = createScriptRequest.Title
	script.Text = createScriptRequest.Text
	script.Owner = user.Id
	script.Create()

	type Response struct {
		Id int64 `json:"id"`
	}

	RespondJSON(200, ctx, Response{Id: script.Id})
}

func UpdateScript(ctx *fasthttp.RequestCtx) {
	type UpdateScriptRequest struct {
		Id    int64  `json:"id"`
		Title string `json:"title"`
		Text  string `json:"text"`
	}

	var updateScriptRequest UpdateScriptRequest
	err := json.Unmarshal(ctx.Request.Body(), &updateScriptRequest)
	if err != nil {
		log.Println(err.Error())
	}

	ctx.SetContentType("text/json")

	userInterface := ctx.UserValue("user")
	user := userInterface.(models.User)

	defer ctx.Request.ConnectionClose()

	script := models.GetScriptById(updateScriptRequest.Id)

	if script.Owner != user.Id {
		RespondJSON(fasthttp.StatusForbidden, ctx, Response{"Forbidden"})
		return
	}

	script.Title = updateScriptRequest.Title
	script.Text = updateScriptRequest.Text

	log.Println(script)

	saveErr := script.Save()

	if saveErr != nil {
		RespondJSON(fasthttp.StatusInternalServerError, ctx, Response{saveErr.Error()})
	} else {
		RespondJSON(fasthttp.StatusOK, ctx, Response{"Saved"})
	}
}

func GetScripts(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/json")

	userInterface := ctx.UserValue("user")
	user := userInterface.(models.User)

	userScripts, _ := models.GetUserScripts(user)

	RespondJSON(200, ctx, userScripts)
}

func RemoveScript(ctx *fasthttp.RequestCtx) {
	type removeScriptRequest struct {
		Id int64 `json:"id"`
	}

	var request removeScriptRequest
	err := json.Unmarshal(ctx.Request.Body(), &request)
	if err != nil {
	}

	ctx.SetContentType("text/json")

	userInterface := ctx.UserValue("user")
	user := userInterface.(models.User)

	defer ctx.Request.ConnectionClose()

	models.RemoveScript(request.Id, user.Id)

	RespondJSON(200, ctx, Response{"Removed"})
}

func GetScriptsAnswers(ctx *fasthttp.RequestCtx) {
	leadId, _ := strconv.ParseInt(string(ctx.QueryArgs().Peek("lead")[:]), 10, 64)

	ctx.SetContentType("text/json")

	userInterface := ctx.UserValue("user")
	user := userInterface.(models.User)

	defer ctx.Request.ConnectionClose()

	answers, getErr := models.GetLeadAnswers(leadId, user.Id)

	if getErr != nil {
		RespondJSON(fasthttp.StatusInternalServerError, ctx, Response{getErr.Error()})
	} else {
		RespondJSON(fasthttp.StatusOK, ctx, answers)
	}
}
