package main

import (
	"fmt"
	"net/http"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	"taskWorker/application"
	"taskWorker/common"
	"taskWorker/repo"
	"taskWorker/server"
	"taskWorker/service1"
)

func main() {
	log.Info("Retrieving configuration")
	cfg, err := common.GetConfig()
	if nil != err {
		panic(fmt.Errorf("cannot get env variables: %w", err))
	}

	app := application.App{
		Config: cfg,
	}

	log.Info("Connecting to local repo")
	r, err := repo.NewRepo(app.Config)
	if nil != err {
		panic(fmt.Errorf("cannot connect to MYSQL server: %w", err))
	}

	app.Repo = r

	app.Service1 = service1.New("http://localhost:3001", nil)

	app.Cron = cron.New()

	if app.Config.WrkEnabled {
		go startRequestWrk(&app)
		go startJobWrk(&app)
	}

	app.Cron.Start()

	log.WithField("Port", app.Config.Port).Info("Starting server")
	panic(http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Port), server.NewServer()))
}
