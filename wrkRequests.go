package main

import (
	"fmt"

	"github.com/TheFranMan/tasker-common/types"
	log "github.com/sirupsen/logrus"

	"worker/application"
)

func startRequestWrk(app *application.App) {
	log.WithField("cron", app.Config.WrkRequestNewCron).Info("Starting request worker")

	app.Cron.AddFunc(app.Config.WrkRequestNewCron, func() {
		log.Debug("Starting new request run")

		err := processNewRequests(app)
		if nil != err {
			log.WithError(err).Error("cannot process new requests")
		}
	})

	app.Cron.AddFunc(app.Config.WrkRequestInProgressCron, func() {
		log.Debug("Starting in progress request run")

		err := processInProgressRequests(app)
		if nil != err {
			log.WithError(err).Error("cannot process in progress requests")
		}
	})
}

func processNewRequests(app *application.App) error {
	requests, err := app.Repo.GetNewRequests()
	if nil != err {
		return fmt.Errorf("cannot retrieve new requests: %w", err)
	}

	for _, request := range requests {
		jobs := []types.Job{}
		l := log.WithFields(log.Fields{
			"token":  request.Token,
			"action": request.Action,
		})

		for _, job := range request.Steps[0].Jobs {
			jobs = append(jobs, types.Job{
				Token: request.Token,
				Name:  job,
				Step:  0,
			})
		}

		err = app.Repo.InsertJobs(jobs)
		if nil != err {
			l.WithError(err).Error("cannot insert jobs")
			continue
		}

		err = app.Repo.MarkRequestInProgress(request.Token)
		if nil != err {
			l.WithError(err).Error("cannot mark request as in progress")
			continue
		}
	}

	return nil
}

func processInProgressRequests(app *application.App) error {
	requests, err := app.Repo.GetInProgressRequests()
	if nil != err {
		return fmt.Errorf("cannot get in progress requests: %w", err)
	}

	for _, request := range requests {
		l := log.WithFields(log.Fields{
			"token":  request.Token,
			"step":   request.Step,
			"action": request.Action,
		})

		jobs, err := app.Repo.GetRequestStepJobs(request.Token, request.Step)
		if nil != err {
			l.WithError(err).Error("cannot get jobs for this step")
			continue
		}

		retryJobs := []types.Job{}
		successCnt := 0

		for _, job := range jobs {
			l = l.WithField("name", job.Name)

			// Job completed successfully, keep a count of these so we can check all of the request jobs completed successfully.
			if types.JobStatusCompleted == types.JobStatus(job.Status) {
				successCnt++
				continue
			}

			// If job is failed, mark the request as failed and stop checking the rest of the jobs
			if types.JobStatusFailed == types.JobStatus(job.Status) {
				err = app.Repo.MarkRequestFailed(job.Token)
				if nil != err {
					l.WithError(err).Error("cannot mark request as failed")
				}

				break
			}

			//	If the job has an error, mark the job for reinsertion.
			if types.JobStatusRetry == types.JobStatus(job.Status) {
				retryJobs = append(retryJobs, job)

				continue
			}
		}

		if 0 < len(retryJobs) {
			err = app.Repo.InsertJobs(retryJobs)
			if nil != err {
				l.WithError(err).Error("cannot reinsert failed jobs")
			}

			continue
		}

		if successCnt != len(jobs) {
			continue
		}

		// Mark the request as completed.
		if request.IsLastStep() {
			err = app.Repo.MarkRequestCompleted(request.Token)
			if nil != err {
				l.WithError(err).Error("cannot mark request as completed")
			}

			continue
		}

		// Insert the request's jobs from it's next step.
		jobs = []types.Job{}
		nextStep := request.Step + 1

		for _, job := range request.Steps[nextStep].Jobs {
			jobs = append(jobs, types.Job{
				Token: request.Token,
				Name:  job,
				Step:  nextStep,
			})
		}

		err = app.Repo.InsertJobs(jobs)
		if nil != err {
			l.WithError(err).Error("cannot insert next step jobs")
			continue
		}

		err = app.Repo.UpdateRequestStep(request.Token)
		if nil != err {
			l.WithError(err).Error("cannot update request step")
			continue
		}
	}

	return nil
}
