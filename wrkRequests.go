package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"taskWorker/application"
	"taskWorker/repo"
)

func startRequestWrk(app *application.App) {
	log.WithField("cron", app.Config.WrkRequestNewCron).Info("Starting request worker")

	app.Cron.AddFunc(app.Config.WrkRequestNewCron, func() {
		log.Debug("Starting new request run")

		err := processNewRequests(app)
		if nil != err {
			log.WithError(err).Error("cannot start processing new requests")
		}
	})

	app.Cron.AddFunc(app.Config.WrkRequestInProgressCron, func() {
		log.Debug("Starting inprogress request run")

		err := processInProgressRequests(app)
		if nil != err {
			log.WithError(err).Error("cannot start processing inprogress requests")
		}
	})
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

		errorJobs := []repo.JobDetails{}
		for _, job := range jobs {
			// If job is failed, mark the request as failed and stop checking the rest of the jobs
			if repo.JobStatusFailed == repo.JobStatus(job.Status) {
				err = app.Repo.MarkRequestFailed(request.Token)
				if nil != err {
					return err
				}

				break
			}

			//	If the job has an error, reinsert the job.
			if repo.JobStatusError == repo.JobStatus(job.Status) {
				errorJobs = append(errorJobs, repo.JobDetails{
					Token: job.Token,
					Name:  job.Name,
					Step:  job.Step,
				})

				continue
			}

		}

		if 0 < len(errorJobs) {
			err = app.Repo.InsertJobs(errorJobs)
			if nil != err {
				return err
			}

			continue
		}

		// Launch the next round of jobs
		if !request.IsLastStep() {
			jobs := []repo.JobDetails{}
			nextStep := request.Step + 1

			for _, job := range request.Steps[nextStep].Jobs {
				jobs = append(jobs, repo.JobDetails{
					Token: request.Token,
					Name:  job,
					Step:  nextStep,
				})
			}

			err = app.Repo.InsertJobs(jobs)
			if nil != err {
				return fmt.Errorf("cannot insert inital jobs: %w", err)
			}

			err = app.Repo.UpdateRequestStep(request.Token)
			if nil != err {
				return err
			}

			continue
		}

		// Mark as completed
		err = app.Repo.MarkRequestCompleted(request.Token)
		if nil != err {
			return err
		}

	}

	return nil
}
