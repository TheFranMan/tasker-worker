package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"taskWorker/application"
	"taskWorker/common"
	repo "taskWorker/repo"
	"taskWorker/server"
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

	if app.Config.StartWorkers {
		go startRequestWrk(&app)
		go startJobWrk(&app)
	}

	log.WithField("Port", app.Config.Port).Info("Starting server")
	panic(http.ListenAndServe(fmt.Sprintf(":%d", app.Config.Port), server.NewServer()))
}
