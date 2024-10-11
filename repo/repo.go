package repo

import (
	"fmt"
	"slices"
	"time"

	"github.com/TheFranMan/tasker-common/types"
	_ "github.com/go-sql-driver/mysql"
	sqlb "github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"

	"worker/common"
)

type Interface interface {
	GetNewRequests() ([]types.Request, error)
	GetInProgressRequests() ([]types.Request, error)
	GetRequest(token string) (*types.Request, error)
	SaveExtra(key string, value any, token string) error
	GetNewJobs() ([]types.Job, error)
	GetRequestStepJobs(token string, step int) ([]types.Job, error)
	MarkRequestFailed(token string) error
	MarkRequestCompleted(token string) error
	InsertJobs(jobs []types.Job) error
	MarkRequestInProgress(token string) error
	UpdateRequestStep(token string) error
	MarkJobNew(id int) error
	MarkJobInprogress(id int) error
	MarkJobCompleted(id int) error
	MarkJobRetry(id int, err error) error
	MarkJobFailed(id int, err error) error
}

type Repo struct {
	db *sqlx.DB
}

func NewRepo(config *common.Config) (*Repo, error) {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DbUser, config.DbPass, config.DbHost, config.DbPort, config.DbName))
	if nil != err {
		return nil, err
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(5 * time.Minute)

	return &Repo{db}, nil
}

func NewRepoWithDb(db *sqlx.DB) *Repo {
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(5 * time.Minute)

	return &Repo{db}
}

func (r *Repo) GetNewRequests() ([]types.Request, error) {
	return r.getRequests(types.RequestStatusNew)
}

func (r *Repo) GetInProgressRequests() ([]types.Request, error) {
	return r.getRequests(types.RequestStatusInProgress)
}

func (r *Repo) GetRequest(token string) (*types.Request, error) {
	var request types.Request
	err := r.db.Get(&request, "SELECT token, request_token, action, params, extras, steps, step, created, completed FROM requests WHERE token = ?", token)
	return &request, err
}

func (r *Repo) GetRequestStepJobs(token string, step int) ([]types.Job, error) {
	var jobs []types.Job
	err := r.db.Select(&jobs, `SELECT j.id, j.name, j.token, j.step, j.error, j.status, j.created, j.completed
		FROM jobs j
		INNER JOIN (
			SELECT name, MAX(created) AS max_created
			FROM jobs
			WHERE token = ? AND step = ?
			GROUP BY name
		) latest_jobs ON j.name = latest_jobs.name AND j.created = latest_jobs.max_created
		WHERE token = ?`, token, step, token)
	if nil != err {
		return nil, fmt.Errorf("cannot retrieve jobs: %w", err)
	}

	return jobs, nil
}

func (r *Repo) InsertJobs(jobs []types.Job) error {
	ib := sqlb.NewInsertBuilder()
	ib.InsertInto("jobs")
	ib.Cols("token", "name", "step", "status", "created")

	for _, job := range jobs {
		ib.Values(job.Token, job.Name, job.Step, 0, sqlb.Raw("NOW()"))
	}

	sql, args := ib.Build()
	_, err := r.db.Exec(sql, args...)
	return err
}

func (r *Repo) SaveExtra(key string, value any, token string) error {
	_, err := r.db.Exec(fmt.Sprintf("UPDATE requests SET extras = JSON_SET(IFNULL(extras, '{}'), '$.%s', ?) WHERE token = ?", key), value, token)
	return err
}

func (r *Repo) MarkRequestFailed(token string) error {
	return r.updateRequestStatus(token, types.RequestStatusFailed)
}

func (r *Repo) MarkRequestInProgress(token string) error {
	return r.updateRequestStatus(token, types.RequestStatusInProgress)
}

func (r *Repo) MarkRequestCompleted(token string) error {
	return r.updateRequestStatus(token, types.RequestStatusCompleted)
}

func (r *Repo) UpdateRequestStep(token string) error {
	_, err := r.db.Exec("UPDATE requests SET step = step + 1 WHERE token = ?", token)
	return err
}

func (r *Repo) GetNewJobs() ([]types.Job, error) {
	return r.getJobs(types.JobStatusNew)
}

func (r *Repo) MarkJobNew(id int) error {
	return r.updateJobStatus(id, types.JobStatusNew)
}

func (r *Repo) MarkJobInprogress(id int) error {
	return r.updateJobStatus(id, types.JobStatusInProgress)
}

func (r *Repo) MarkJobCompleted(id int) error {
	return r.updateJobStatus(id, types.JobStatusCompleted)
}

func (r *Repo) MarkJobRetry(id int, err error) error {
	return r.markJobWithError(id, err, types.JobStatusRetry)
}

func (r *Repo) MarkJobFailed(id int, err error) error {
	return r.markJobWithError(id, err, types.JobStatusFailed)
}

func (r *Repo) markJobWithError(id int, err error, status types.JobStatus) error {
	_, err = r.db.Exec("Update jobs SET status = ?, error = ? WHERE id = ?", status, err.Error(), id)
	if nil != err {
		return err
	}

	return nil
}

func (r *Repo) getRequests(status types.RequestStatus) ([]types.Request, error) {
	var requests []types.Request
	err := r.db.Select(&requests, "SELECT token, request_token, action, params, extras, steps, step, status, created, completed FROM requests WHERE status = ?", int(status))
	if nil != err {
		return nil, err
	}

	return requests, nil
}

func (r *Repo) updateRequestStatus(token string, status types.RequestStatus) error {
	sets := []string{}

	ub := sqlb.NewUpdateBuilder()
	ub.Update("requests")
	sets = append(sets, ub.Assign("status", status))

	if slices.Contains([]types.RequestStatus{types.RequestStatusCompleted, types.RequestStatusFailed}, status) {
		sets = append(sets, "completed = NOW()")
	}

	ub.Set(sets...)
	ub.Where(ub.Equal("token", token))

	sql, args := ub.Build()
	_, err := r.db.Exec(sql, args...)
	return err
}

func (r *Repo) updateJobStatus(id int, status types.JobStatus) error {
	ub := sqlb.NewUpdateBuilder()
	ub.Update("jobs")
	ub.Where(ub.Equal("id", id))

	params := []string{ub.Assign("status", status)}
	if types.JobStatusNew != status {
		params = append(params, "completed = NOW()")
	}

	ub.Set(params...)

	sql, args := ub.Build()

	_, err := r.db.Exec(sql, args...)
	return err
}

func (r *Repo) getJobs(status types.JobStatus) ([]types.Job, error) {
	var jobs []types.Job
	err := r.db.Select(&jobs, "SELECT id, name, token FROM jobs WHERE status = ?", status)
	if nil != err {
		return nil, fmt.Errorf("cannot select new jobs: %w", err)
	}

	return jobs, nil
}
