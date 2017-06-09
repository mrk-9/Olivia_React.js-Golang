package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/fasthttp-contrib/websocket"
	"github.com/valyala/fasthttp"
	"olivia/models"
	"strings"
	"log"
)

type Response struct {
	Message string `json:"message"`
}

var Router *fasthttprouter.Router

var MiddlewareFuncs []func(*fasthttp.RequestCtx) (*fasthttp.RequestCtx, error)

var upgrader = websocket.New(WS)

func Init() {
	Router.GET("/ws", func(ctx *fasthttp.RequestCtx) {
		upgradeErr := upgrader.Upgrade(ctx)
		if upgradeErr != nil {
			log.Printf("Cannot upgrade: %s", upgradeErr.Error())
		} else {
			return
		}
	})

	Router.POST("/auth/signup", Register)
	Router.POST("/auth/signin", Login)
	Router.GET("/auth/getprofile", GetProfile)
	Router.POST("/auth/validatephone", ValidatePhone)

	Router.GET("/emails/get", GetEmails)

	Router.GET("/leads/get", GetLeads)
	Router.POST("/leads/setownership", SetOwnership)

	Router.POST("/sms", ReceiveSMS)
	Router.POST("/sms/send", SendSMS)
	Router.GET("/sms/chat", GetChat)

	Router.POST("/scripts/create", CreateScript)
	Router.GET("/scripts/get", GetScripts)
	Router.POST("/scripts/remove", RemoveScript)
	Router.POST("/scripts/update", UpdateScript)
	Router.GET("/scripts/answers", GetScriptsAnswers)


	Router.POST("/support/auth/signin", SupportLogin)
	Router.GET("/support/auth/getprofile", SupportGetProfile)
	Router.GET("/support/leads/get", GetSupportLeads)

	Router.GET("/support/scriptsandanswers", GetScriptsAndAnswers)

	Router.GET("/support/sms/chat", GetSupportChat)
	Router.POST("/support/sms/send", SendSupportSMS)

	Router.POST("/support/scripts/answer", AnswerScript)

	Router.POST("/support/leads/setassistance", SetAssistance)

	MiddlewareFuncs = append(MiddlewareFuncs, func(ctx *fasthttp.RequestCtx) (*fasthttp.RequestCtx, error) {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Response.Header.Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		return ctx, nil
	})

	// Get user by Auth token if available
	MiddlewareFuncs = append(MiddlewareFuncs, func(ctx *fasthttp.RequestCtx) (*fasthttp.RequestCtx, error) {
		if string(ctx.Method()) == "OPTIONS" {
			return ctx, nil
		}

		// Ignore URIs that does not require authentication
		ignoreURIs := [...]string{"/support/auth/signin", "/auth/signup", "/auth/signin", "/", "/sms", "/auth/validatephone"}
		for _, uri := range ignoreURIs {
			if uri == string(ctx.URI().Path()) {
				return ctx, nil
			}
		}

		authorization := ctx.Request.Header.Peek("Authorization")
		authSplit := strings.Split(string(authorization), " ")

		if len(authSplit) != 2 {
			RespondJSON(400, ctx, Response{"Authentication token is invalid"})
			return ctx, errors.New("Authentication token is invalid")
		}

		key, token := authSplit[0], authSplit[1]

		if token != "" {
			if (key == "Bearer") {
				user := models.User{}
				user.GetByToken(token)

				// Fail if user not found
				if user.Id != 0 {
					ctx.SetUserValue("user", user)
				} else {
					RespondJSON(400, ctx, Response{"Authentication token is invalid"})
					return ctx, errors.New("Authentication token is invalid")
				}
			} else if (key == "SupportBearer") {
				user := models.SupportUser{}
				user.GetByToken(token)

				// Fail if user not found
				if user.Id != 0 {
					ctx.SetUserValue("user", user)
				} else {
					RespondJSON(400, ctx, Response{"Authentication token is invalid"})
					return ctx, errors.New("Authentication token is invalid")
				}
			}
		}

		return ctx, nil
	})
}

func HandleErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func RespondJSON(status int, ctx *fasthttp.RequestCtx, response interface{}) {
	ctx.SetStatusCode(status)
	jsonBytes, err := json.Marshal(response)
	HandleErr(err)

	jsonString := string(jsonBytes)

	fmt.Fprint(ctx, jsonString)

	ctx.SetConnectionClose()
}
