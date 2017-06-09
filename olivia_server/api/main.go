package main

import (
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"flag"
	"olivia/models"
	"olivia/routes"
	"strconv"
	"fmt"
	"config"
	"olivia/zendesk"
)

var DB *sql.DB

func HandleErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

type Middleware struct {
	Router    fasthttprouter.Router
	Functions []func(*fasthttp.RequestCtx) (*fasthttp.RequestCtx, error)
}

func (mw Middleware) Handler(ctx *fasthttp.RequestCtx) {
	var err error

	if string(ctx.Path()) != "/ws" {
		for _, f := range mw.Functions {
			ctx, err = f(ctx)
			if err != nil {
				return
			}
		}
	}

	mw.Router.Handler(ctx)
}

func main() {
	var err error
	DB, err = sql.Open("mysql", config.MysqlHost)
	HandleErr(err)

	models.DB = DB

	DB.Exec("SET NAMES `utf8`")
	DB.Exec("SET CHARACTER SET utf8")

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	zendesk.SetCredentials(config.ZDSubdomain, config.ZDUsername, config.ZDPassword)

	models.CreateUsersTable()
	models.CreateTokensTable()
	models.CreateEmailsTable()
	models.CreateLeadsTable()
	models.CreatePropertyTable()
	models.CreatePhoneNumbersTable()
	models.CreateMessagesTable()
	models.CreateScriptsTable()
	models.CreateScriptAnswersTable()
	models.CreateBadEmailsTable()

	models.CreateSupportUsersTable()
	models.CreateSupportTokensTable()

	router := fasthttprouter.New()
	routes.Router = router
	routes.Init()

	port := flag.Int("port", 8081, "Port to listen, 8081 is default")
	flag.Parse()

	var mw = Middleware{*router, routes.MiddlewareFuncs}

	fmt.Println("Port used: " + strconv.Itoa(*port))
	log.Fatal(fasthttp.ListenAndServe(":"+strconv.Itoa(*port), mw.Handler))
}
