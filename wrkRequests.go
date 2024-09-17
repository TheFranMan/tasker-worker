package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"taskWorker/application"
	"taskWorker/repo"
)

func startRequestWrk(app *application.App) {
	log.Info("Starting request worker")

	err := processNewRequests(app)
	if nil != err {
		log.WithError(err).Error("cannot start processing new requests")
	}

	err = processInProgressRequests(app)
	if nil != err {
		log.WithError(err).Error("cannot start processing inprogress requests")
	}
}

func processNewRequests(app *application.App) error {
	requests, err := app.Repo.GetNewRequests()
	if nil != err {
		return err
	}

	for _, request := range requests {
		// fmt.Printf("%+v\n", request)
		jobs := []repo.JobDetails{}

		for _, job := range request.Steps[0].Jobs {
			jobs = append(jobs, repo.JobDetails{
				Token: request.Token,
				Name:  job,
				Step:  0,
			})
		}

		err = app.Repo.InsertJobs(jobs)
		if nil != err {
			return fmt.Errorf("cannot insert inital jobs: %w", err)
		}

		err = app.Repo.MarkRequestInProgress(request.Token)
		if nil != err {
			return fmt.Errorf("cannot mark request as in progress: %w", err)
		}
	}

	return nil
}

func processInProgressRequests(app *application.App) error {
	requests, err := app.Repo.GetInProgressRequests()
	if nil != err {
		return err
	}

	for _, request := range requests {
		jobs, err := app.Repo.GetRequestStepJobs(request.Token, request.Step)
		if nil != err {
			return err
		}
		_ = jobs
		// for _, job := range jobs {

		//}
	}

	return nil
}
