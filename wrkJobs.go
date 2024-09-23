package main

import (
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

	callbacks["service1DeleteUser"] = func(app *application.App, request repo.Request, id int) error {
		return cb.Service1DeleteUser(app, request, id)
	}
}

func startJobWrk(app *application.App) {
	log.WithField("cron", app.Config.WrkJobCron).Info("Starting job worker")

	app.Cron.AddFunc(app.Config.WrkJobCron, func() {
		log.Debug("Starting Job worker run")

		err := processNewJobs(app)
		if nil != err {
			log.WithError(err).Error("cannot process new jobs")
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
			log.WithError(err).Error("cannot retrieve a jobs request: %w", err)
			continue
		}

		l := log.WithFields(log.Fields{
			"id":       job.ID,
			"token":    job.Token,
			"step":     job.Step,
			"callback": job.Name,
		})

		if nil == request {
			l.Warn("cannot retrieve request based on a job")
			continue
		}

		if _, exists := callbacks[job.Name]; !exists {
			l.Warn("unknown job name")
			continue
		}

		tryErr := func() error {
			// Update job status
			err := app.Repo.MarkJobInprogress(job.ID)
			if nil != err {
				return fmt.Errorf("cannot update job status to inprogress: %w", err)
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
				l.WithError(fmt.Errorf("reset job error: %w, original error: %w", tryErr, err)).Errorf("cannot reset job as new after orginal error")
				continue
			}

			log.WithError(tryErr).Error("cannot process new job")
			continue
		}
	}

	return nil
}
