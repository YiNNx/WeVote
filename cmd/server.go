package main

import (
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"

	"github.com/YiNNx/WeVote/internal/config"
	"github.com/YiNNx/WeVote/internal/gqlgen"
	"github.com/YiNNx/WeVote/internal/jobs"
	"github.com/YiNNx/WeVote/internal/models"
	"github.com/YiNNx/WeVote/internal/schema"
	"github.com/YiNNx/WeVote/internal/utils/ticket"
	"github.com/YiNNx/WeVote/pkg/captcha"
	"github.com/YiNNx/WeVote/pkg/log"
)

func main() {
	confPath := os.Getenv("CONF_PATH")
	conf := config.Init(confPath)

	log.InitLogger(conf.Log.Path, conf.Server.DebugMode)
	ticket.InitGenerator(conf.Ticket.Secret, conf.Ticket.Spec)
	jobs.InitJobs(conf.Ticket.Spec)
	captcha.InitReChaptcha(conf.Captcha.RecaptchaSecret)
	models.InitDataBaseConnections(
		conf.Postgres.Host,
		conf.Postgres.Port,
		conf.Postgres.User,
		conf.Postgres.Password,
		conf.Postgres.Dbname,
		conf.Redis.Addrs,
	)

	http.Handle("/", handler.NewDefaultServer(
		gqlgen.NewExecutableSchema(
			gqlgen.Config{Resolvers: &schema.Resolver{}},
		)))
	log.Logger.Fatal(http.ListenAndServe(conf.Server.Addr, nil))
}
