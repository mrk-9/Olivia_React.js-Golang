package routes

import (
	"github.com/valyala/fasthttp"
	"olivia/models"
)

func GetEmails(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/json")

	userInterface := ctx.UserValue("user")
	user := userInterface.(models.User)

	emails := models.GetEmailsByUserId(user.Id)

	RespondJSON(200, ctx, emails)
}
