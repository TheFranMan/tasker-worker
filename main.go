package main

import (
	"fmt"
	"net/http"
	"taskWorker/common"
	"taskWorker/server"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Retrieving configuration")
	cfg, err := common.GetConfig()
	if nil != err {
		panic(fmt.Errorf("cannot get env variables: %w", err))
	}

	log.WithField("Port", cfg.Port).Info("Starting server")
	panic(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), server.NewServer()))
}
