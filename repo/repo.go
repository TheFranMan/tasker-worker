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

type requestStatus int
type jobStatus int

var (
	requestStatusNew requestStatus = 0

	jobStatusNew        jobStatus = 0
	jobStatusInProgress jobStatus = 1
)

type Interface interface {
	GetNewRequests() ([]Request, error)
	GetNewJobs() ([]Job, error)
	InsertJobs(jobDetails []JobDetails) error
	GetRequest(token string) (*Request, error)
	MarkJobsInprogress(id int) error
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
	ID        int          `db:"id"`
	Name      string       `db:"name"`
	Token     string       `db:"token"`
	Step      int          `db:"step"`
	Status    int          `db:"status"`
	Created   time.Time    `db:"created"`
	Completed sql.NullTime `db:"complete"`
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

type Step struct {
	Name string   `json:"name"`
	Jobs []string `json:"jobs"`
}

type Steps []Step

func (s *Steps) Scan(value interface{}) error {
	if nil == value {
		s = &Steps{}
		return nil
	}

	return json.Unmarshal(value.([]byte), s)
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

func (r *Repo) GetNewRequests() ([]Request, error) {
	return r.getRequests(requestStatusNew)
}

func (r *Repo) getRequests(status requestStatus) ([]Request, error) {
	var requests []Request
	err := r.db.Select(&requests, "SELECT token, request_token, action, params, extras, steps, step, status, created, completed FROM requests WHERE status = ?", int(status))
	if nil != err {
		return nil, err
	}

	return requests, nil
}

type JobDetails struct {
	Token string
	Name  string
	Step  int
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

func (r *Repo) GetNewJobs() ([]Job, error) {
	return r.getJobs(jobStatusNew)
}

func (r *Repo) getJobs(status jobStatus) ([]Job, error) {
	var jobs []Job
	err := r.db.Select(&jobs, "SELECT id, name, token FROM jobs WHERE status = ?", status)
	if nil != err {
		return nil, fmt.Errorf("cannot select new jobs: %w", err)
	}

	return jobs, nil
}

func (r *Repo) GetRequest(token string) (*Request, error) {
	var request Request
	err := r.db.Get(&request, "SELECT * FROM requests WHERE token = ?", token)
	return &request, err
}

func (r *Repo) MarkJobsInprogress(id int) error {
	return r.updateJobStatus(id, jobStatusInProgress)
}

func (r *Repo) updateJobStatus(id int, status jobStatus) error {
	_, err := r.db.Exec("UPDATE jobs SET status = ? WHERE id = ?", status, id)
	return err
}
