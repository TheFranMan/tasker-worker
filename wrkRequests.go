package main

import (
	"fmt"
	"time"

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

	time.Sleep(time.Second)

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

			//	If all jobs completed successfully, send the request next step jobs. If on the final step, mark the job as completed.
		}

		fmt.Printf("%+v\n", errorJobs)

		if 0 < len(errorJobs) {
			err = app.Repo.InsertJobs(errorJobs)
			if nil != err {
				return err
			}

			continue
		}
	}

	return nil
}
