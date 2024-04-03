package main

import (
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"

	"github.com/YiNNx/WeVote/internal/common/log"
	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/internal/models"
	"github.com/YiNNx/WeVote/internal/schema"
	"github.com/YiNNx/WeVote/internal/schema/gqlgen"
	"github.com/YiNNx/WeVote/pkg/captcha"
)

func main() {
	confPath := os.Getenv("CONF_PATH")
	err := config.Init(confPath)
	if err != nil {
		log.Logger.Fatal(err)
	}

	log.InitLogger(config.C.Log.Path, config.C.Server.DebugMode)
	captcha.NewClient(config.C.Captcha.RecaptchaSecret)
	models.InitConnections(config.C.Postgres.DSN, config.C.Redis.Addrs)

	http.Handle("/", handler.NewDefaultServer(
		gqlgen.NewExecutableSchema(
			gqlgen.Config{Resolvers: &schema.Resolver{}},
		)))
	log.Logger.Fatal(http.ListenAndServe(config.C.Server.Addr, nil))
}
