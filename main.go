package main

import (
	"fmt"
	"net/http"
	"taskWorker/common"
	"taskWorker/server"

	log "github.com/sirupsen/logrus"
)

func main() {
	env, err := common.GetEnv()
	if nil != err {
		panic(fmt.Errorf("cannot get env variables: %w", err))
	}

	log.WithField("Port", env.Port).Info("Starting server")
	panic(http.ListenAndServe(fmt.Sprintf(":%d", env.Port), server.NewServer()))

	fmt.Printf("%+v\n", env)
}
