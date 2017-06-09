package routes

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"olivia/models"
)

func SupportLogin(ctx *fasthttp.RequestCtx) {
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

	var user models.SupportUser
	user.GetByCredentials(loginRequest.Email, loginRequest.Password)

	if user.Id == 0 {
		RespondJSON(400, ctx, Response{"User not found"})
		return
	}

	var token models.SupportToken
	token.GetTokenForUser(user)

	if token.Id == 0 {
		token.Create(user)
	}

	RespondJSON(200, ctx, token)
}

func SupportGetProfile(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/json")

	userInterface := ctx.UserValue("user")
	user := userInterface.(models.SupportUser)

	RespondJSON(200, ctx, user)
}