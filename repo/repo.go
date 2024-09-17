package repo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	sqlb "github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"

	"taskWorker/common"
)

type Interface interface {
	GetNewRequests() ([]Request, error)
	GetInProgressRequests() ([]Request, error)
	GetRequest(token string) (*Request, error)
	SaveExtra(key, value, token string) error
	GetNewJobs() ([]Job, error)
	GetRequestStepJobs(token string, step int) ([]Job, error)
	MarkRequestFailed(token string) error
	InsertJobs(jobDetails []JobDetails) error
	MarkRequestInProgress(token string) error
	MarkJobNew(id int) error
	MarkJobInprogress(id int) error
	MarkJobCompleted(id int) error
	MarkJobFailed(id int) error
}

type RequestStatus int
type JobStatus int

type JobDetails struct {
	Token string
	Name  string
	Step  int
}

type Request struct {
	Token        string         `db:"token"`
	RequestToken string         `db:"request_token"`
	Action       string         `db:"action"`
	Params       Params         `db:"params"`
	Extras       sql.NullString `db:"extras"`
	Steps        Steps          `db:"steps"`
	Step         int            `db:"step"`
	Status       int            `db:"status"`
	Created      time.Time      `db:"created"`
	Completed    sql.NullTime   `db:"completed"`
}

type Job struct {
	ID        int            `db:"id"`
	Name      string         `db:"name"`
	Token     string         `db:"token"`
	Step      int            `db:"step"`
	Error     sql.NullString `db:"error"`
	Status    int            `db:"status"`
	Created   time.Time      `db:"created"`
	Completed sql.NullTime   `db:"completed"`
}

type Params struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (p *Params) Scan(value interface{}) error {
	if nil == value {
		p = &Params{}
		return nil
	}

	return json.Unmarshal(value.([]byte), p)
}

type Steps []Step
type Step struct {
	Name string   `json:"name"`
	Jobs []string `json:"jobs"`
}

func (s *Steps) Scan(value interface{}) error {
	if nil == value {
		s = &Steps{}
		return nil
	}

	return json.Unmarshal(value.([]byte), s)
}

var (
	RequestStatusNew        RequestStatus = 0
	RequestStatusInProgress RequestStatus = 1
	RequestStatusFailed     RequestStatus = 3

	JobStatusNew        JobStatus = 0
	JobStatusInProgress JobStatus = 1
	JobStatusCompleted  JobStatus = 2
	JobStatusFailed     JobStatus = 4
)

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

func (r *Repo) GetNewRequests() ([]Request, error) {
	return r.getRequests(RequestStatusNew)
}

func (r *Repo) GetInProgressRequests() ([]Request, error) {
	return r.getRequests(RequestStatusInProgress)
}

func (r *Repo) GetRequest(token string) (*Request, error) {
	var request Request
	err := r.db.Get(&request, "SELECT token, request_token, action, params, extras, steps, step, created, completed FROM requests WHERE token = ?", token)
	return &request, err
}

func (r *Repo) GetRequestStepJobs(token string, step int) ([]Job, error) {
	var jobs []Job
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

func (r *Repo) InsertJobs(jobDetails []JobDetails) error {
	ib := sqlb.NewInsertBuilder()
	ib.InsertInto("jobs")
	ib.Cols("token", "name", "step", "status", "created")
	for _, jobDetail := range jobDetails {
		ib.Values(jobDetail.Token, jobDetail.Name, jobDetail.Step, 0, sqlb.Raw("NOW()"))
	}

	sql, args := ib.Build()
	_, err := r.db.Exec(sql, args...)
	return err
}

func (r *Repo) SaveExtra(key, value, token string) error {
	_, err := r.db.Exec(fmt.Sprintf("UPDATE requests SET extras = JSON_SET(extras, '$.%s', ?) WHERE token = ?", key), value, token)
	return err
}
func (r *Repo) MarkRequestFailed(token string) error {
	return r.updateRequestStatus(token, RequestStatusFailed)
}
func (r *Repo) MarkRequestInProgress(token string) error {
	return r.updateRequestStatus(token, RequestStatusInProgress)
}

func (r *Repo) GetNewJobs() ([]Job, error) {
	return r.getJobs(JobStatusNew)
}

func (r *Repo) MarkJobNew(id int) error {
	return r.updateJobStatus(id, JobStatusNew)
}

func (r *Repo) MarkJobInprogress(id int) error {
	return r.updateJobStatus(id, JobStatusInProgress)
}

func (r *Repo) MarkJobCompleted(id int) error {
	return r.updateJobStatus(id, JobStatusCompleted)
}

func (r *Repo) MarkJobFailed(id int) error {
	return r.updateJobStatus(id, JobStatusFailed)
}

func (r *Repo) getRequests(status RequestStatus) ([]Request, error) {
	var requests []Request
	err := r.db.Select(&requests, "SELECT token, request_token, action, params, extras, steps, step, status, created, completed FROM requests WHERE status = ?", int(status))
	if nil != err {
		return nil, err
	}

	return requests, nil
}

func (r *Repo) updateRequestStatus(token string, status RequestStatus) error {
	_, err := r.db.Exec("UPDATE requests SET status = ? WHERE token = ?", status, token)
	return err
}

func (r *Repo) updateJobStatus(id int, status JobStatus) error {
	ub := sqlb.NewUpdateBuilder()
	ub.Update("jobs")
	ub.Set(
		ub.Assign("status", status),
		"completed = NOW()",
	)
	ub.Where(ub.Equal("id", id))

	sql, args := ub.Build()

	_, err := r.db.Exec(sql, args...)
	return err
}

func (r *Repo) getJobs(status JobStatus) ([]Job, error) {
	var jobs []Job
	err := r.db.Select(&jobs, "SELECT id, name, token FROM jobs WHERE status = ?", status)
	if nil != err {
		return nil, fmt.Errorf("cannot select new jobs: %w", err)
	}

	return jobs, nil
}
