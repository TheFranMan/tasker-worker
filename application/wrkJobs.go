package application

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

func StartJobWrk(app *App) {
	log.Info("Starting job worker")

	err := processNewJobs(app)
	if nil != err {
		log.WithError(err).Error("cannot start processing new jobs")
	}
}

func processNewJobs(app *App) error {
	time.Sleep(time.Second)
	jobs, err := app.Repo.GetNewJobs()
	if nil != err {
		return fmt.Errorf("cannot retrieve new jobs: %w", err)
	}

	fmt.Printf("jobs %+v\n", jobs)
	return nil
}
