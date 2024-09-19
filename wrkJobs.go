package main

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"taskWorker/application"
	cb "taskWorker/callbacks"
	"taskWorker/repo"
)

var callbacks = map[string]func(*application.App, repo.Request, int) error{}

func init() {
	callbacks["service1GetUser"] = func(app *application.App, request repo.Request, id int) error {
		return cb.Service1GetUser(app, request, id)
	}
}

func startJobWrk(app *application.App) {
	log.WithField("cron", app.Config.WrkJobCron).Info("Starting job worker")

	app.Cron.AddFunc(app.Config.WrkJobCron, func() {
		log.Debug("Starting Job worker run")

		err := processNewJobs(app)
		if nil != err {
			log.WithError(err).Error("cannot start processing new jobs")
		}
	})
}

func processNewJobs(app *application.App) error {
	jobs, err := app.Repo.GetNewJobs()
	if nil != err {
		return fmt.Errorf("cannot retrieve new jobs: %w", err)
	}

	for _, job := range jobs {
		// Retrieve parent request
		request, err := app.Repo.GetRequest(job.Token)
		if nil != err {
			return fmt.Errorf("cannot retrieve a jobs request: %w", err)
		}

		if nil == request {
			return errors.New("empty request found for job")
		}

		l := log.WithFields(log.Fields{
			"token":    request.Token,
			"callback": job.Name,
		})

		// Call callback
		if _, exists := callbacks[job.Name]; !exists {
			l.Warn("unknown job name")
			continue
		}

		tryErr := func() error {
			// Update job status
			err := app.Repo.MarkJobInprogress(job.ID)
			if nil != err {
				return fmt.Errorf("cannot update job %d status to inprogress: %w", job.ID, err)
			}

			err = callbacks[job.Name](app, *request, job.ID)
			if nil != err {
				return fmt.Errorf("cannot process job: %w", err)
			}
			return nil
		}()

		if nil != tryErr {
			err = app.Repo.MarkJobNew(job.ID)
			if nil != err {
				return fmt.Errorf("cannot reset job %d as new (error: %w) after orginal error: %w", job.ID, err, tryErr)
			}

			return fmt.Errorf("cannot process new job %d: %w", job.ID, tryErr)
		}
	}

	return nil
}
