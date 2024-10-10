package main

import (
	"fmt"
	"net/http"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	"worker/application"
	"worker/common"
	"worker/repo"
	"worker/server"
	"worker/service1"
	"worker/service2"
	"worker/service3"
)

func main() {
	log.Info("Retrieving configuration")
	cfg, err := common.GetConfig()
	if nil != err {
		log.WithError(err).Panic("cannot get enviroment variables")
	}

	app := application.App{
		Config: cfg,
	}

	log.Info("Connecting to local repo")
	r, err := repo.NewRepo(app.Config)
	if nil != err {
		log.WithError(err).Panic("cannot connect to MYSQL server")
	}

	app.Repo = r

	app.Service1 = service1.New(app.Config.URLService1.String(), nil)
	app.Service2 = service2.New(app.Config.URLService2.String(), nil)
	app.Service3 = service3.New(app.Config.URLService3.String(), nil)

	app.Cron = cron.New()

	if app.Config.WrkEnabled {
		go startRequestWrk(&app)
		go startJobWrk(&app)
	}

	app.Cron.Start()

	log.WithField("Port", app.Config.Port).Info("Starting server")
	err = http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Port), server.NewServer())
	if nil != err {
		log.WithError(err).Panic("cannot start HTTP server")
	}
}
