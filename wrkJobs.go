package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"worker/application"
	cb "worker/callbacks"
	"worker/repo"
)

var callbacks = map[string]func(*application.App, repo.Request, int) error{}

func init() {
	callbacks["service1GetUser"] = func(app *application.App, request repo.Request, id int) error {
		return cb.Service1GetUser(app, request, id)
	}

	callbacks["service1DeleteUser"] = func(app *application.App, request repo.Request, id int) error {
		return cb.Service1DeleteUser(app, request, id)
	}

	callbacks["service2DeleteUser"] = func(app *application.App, request repo.Request, id int) error {
		return cb.Service2DeleteUser(app, request, id)
	}

	callbacks["service3DeleteUser"] = func(app *application.App, request repo.Request, id int) error {
		return cb.Service3DeleteUser(app, request, id)
	}

	callbacks["service1UpdateUser"] = func(app *application.App, request repo.Request, id int) error {
		return cb.Service1UpdateUser(app, request, id)
	}

	callbacks["service2UpdateUser"] = func(app *application.App, request repo.Request, id int) error {
		return cb.Service2UpdateUser(app, request, id)
	}

	callbacks["service3UpdateUser"] = func(app *application.App, request repo.Request, id int) error {
		return cb.Service3UpdateUser(app, request, id)
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
		l := log.WithFields(log.Fields{
			"id":    job.ID,
			"token": job.Token,
			"step":  job.Step,
			"name":  job.Name,
		})

		// Retrieve parent request
		request, err := app.Repo.GetRequest(job.Token)
		if nil != err {
			l.WithError(err).Error("cannot retrieve a job's request")
			continue
		}

		if nil == request {
			l.Warn("cannot retrieve request based on a job")
			continue
		}

		if _, exists := callbacks[job.Name]; !exists {
			l.Warn("unknown job name")
			continue
		}

		err = app.Repo.MarkJobInprogress(job.ID)
		if nil != err {
			l.Error("cannot update job status to in-progress")
			continue
		}

		err = callbacks[job.Name](app, *request, job.ID)
		if nil != err {
			l.WithError(err).Error("cannot process new job")
			continue
		}
	}

	return nil
}
