//go:build intergration

package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/TheFranMan/tasker-common/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/suite"

	"worker/application"
	"worker/repo"
)

type Suite struct {
	suite.Suite
	db       *sqlx.DB
	pool     *dockertest.Pool
	resource *dockertest.Resource
	repo     repo.Interface
}

func TestRun(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	var err error

	s.pool, err = dockertest.NewPool("")
	if nil != err {
		s.FailNowf(err.Error(), "cannot create a new dockertest pool")
	}

	err = s.pool.Client.Ping()
	if nil != err {
		s.FailNowf(err.Error(), "cannot ping dockertest client")
	}

	s.resource, err = s.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "8.0",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=secret",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if nil != err {
		s.FailNowf(err.Error(), "cannot create dockertest mysql s.resource")
	}

	s.resource.Expire(60 * 5)
	mysqlPort := s.resource.GetPort("3306/tcp")

	err = s.pool.Retry(func() error {
		var err error
		s.db, err = sqlx.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql?parseTime=true", mysqlPort))
		if err != nil {
			return err
		}

		return s.db.Ping()
	})
	if nil != err {
		s.FailNowf(err.Error(), "cannot open dockertest mysql connection")
	}

	// Migrations
	driver, err := mysql.WithInstance(s.db.DB, &mysql.Config{})
	if nil != err {
		s.FailNowf(err.Error(), "cannot create migration driver:")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if nil != err {
		s.FailNowf(err.Error(), "cannot create new migration instance")
	}

	err = m.Up()
	if nil != err {
		s.FailNowf(err.Error(), "cannot run mysql migrations")
	}

	// Repo
	s.repo = repo.NewRepoWithDb(s.db)
}

func (s *Suite) TearDownSuite() {
	err := s.pool.Purge(s.resource)
	if nil != err {
		s.FailNowf(err.Error(), "cannot purge dockertest mysql resource")
	}
}

func (s *Suite) AfterTest() {
	s.importFile("truncate.sql")
}

func (s *Suite) importFile(filename string) {
	b, err := os.ReadFile("./repo/testdata/" + filename)
	if nil != err {
		s.FailNowf(err.Error(), "cannot open SQL file: %s", filename)
	}

	statements := strings.Split(strings.TrimSpace(string(b)), ";")

	for _, statement := range statements {
		if 0 == len(statement) {
			continue
		}

		_, err := s.db.Exec(statement + ";")
		if nil != err {
			s.FailNowf(err.Error(), "cannot run SQL statement: %s", statement)
		}
	}
}

func (s *Suite) Test_if_the_request_is_new_insert_the_first_step_jobs_and_update_the_request_status_as_in_progress() {
	s.importFile("delete.sql")

	err := processNewRequests(&application.App{
		Repo: s.repo,
	})
	s.Require().Nil(err)

	var requests []repo.Request
	err = s.db.Select(&requests, "SELECT token, action, step, status FROM requests WHERE status = 1")
	s.Require().Nil(err)

	s.Require().Len(requests, 2)
	s.Require().ElementsMatch([]repo.Request{
		{
			Token:  "2b482d15-6c02-4e7f-bae3-0a8fe1dfb301",
			Action: string(types.ActionDelete),
			Step:   0,
			Status: int(repo.RequestStatusInProgress),
		},
		{
			Token:  "482a2d88-d38a-4509-ac94-beadff53c053",
			Action: string(types.ActionDelete),
			Step:   0,
			Status: int(repo.RequestStatusInProgress),
		},
	}, requests)

	var jobs []repo.Job
	err = s.db.Select(&jobs, `SELECT name, token, step, error, status FROM jobs WHERE token IN ("2b482d15-6c02-4e7f-bae3-0a8fe1dfb301", "482a2d88-d38a-4509-ac94-beadff53c053")`)
	s.Require().Nil(err)

	s.Require().Len(jobs, 2)
	s.Require().ElementsMatch([]repo.Job{
		{
			Name:   "service1GetUser",
			Token:  "2b482d15-6c02-4e7f-bae3-0a8fe1dfb301",
			Step:   0,
			Error:  sql.NullString{},
			Status: 0,
		},
		{
			Name:   "service1GetUser",
			Token:  "482a2d88-d38a-4509-ac94-beadff53c053",
			Step:   0,
			Error:  sql.NullString{},
			Status: 0,
		},
	}, jobs)
}

func (s *Suite) Test_if_the_request_has_an_errored_job_reinsert_the_job_and_do_not_update_the_request_status() {
	s.importFile("delete_with_error_jobs.sql")

	err := processInProgressRequests(&application.App{
		Repo: s.repo,
	})
	s.Require().Nil(err)

	var requests []repo.Request
	err = s.db.Select(&requests, "SELECT token, action, step, status FROM requests WHERE status = 1")
	s.Require().Nil(err)

	s.Require().Len(requests, 1)
	s.Require().ElementsMatch([]repo.Request{
		{
			Token:  "89858b95-21bd-47e3-a03e-9069a7440188",
			Action: string(types.ActionDelete),
			Step:   0,
			Status: int(repo.RequestStatusInProgress),
		},
	}, requests)

	var jobs []repo.Job
	err = s.db.Select(&jobs, `SELECT name, token, step, error, status FROM jobs WHERE token = "89858b95-21bd-47e3-a03e-9069a7440188"`)
	s.Require().Nil(err)

	s.Require().Len(jobs, 2)
	s.Require().ElementsMatch([]repo.Job{
		{
			Name:  "service1GetUser",
			Token: "89858b95-21bd-47e3-a03e-9069a7440188",
			Step:  0,
			Error: sql.NullString{
				Valid:  true,
				String: "test error",
			},
			Status: 4,
		},
		{
			Name:   "service1GetUser",
			Token:  "89858b95-21bd-47e3-a03e-9069a7440188",
			Step:   0,
			Error:  sql.NullString{},
			Status: 0,
		},
	}, jobs)
}

func (s *Suite) Test_if_the_request_has_an_failed_job_exit_and_mark_the_request_as_failed() {
	s.importFile("delete_with_failed_job.sql")

	err := processInProgressRequests(&application.App{
		Repo: s.repo,
	})
	s.Require().Nil(err)

	var requests []repo.Request
	err = s.db.Select(&requests, "SELECT token, action, step, status FROM requests WHERE status = 3")
	s.Require().Nil(err)

	s.Require().Len(requests, 1)
	s.Require().ElementsMatch([]repo.Request{
		{
			Token:  "c639b525-1ab1-44e6-bde3-96238cf13f2f",
			Action: string(types.ActionDelete),
			Step:   0,
			Status: int(repo.RequestStatusFailed),
		},
	}, requests)

	var jobs []repo.Job
	err = s.db.Select(&jobs, `SELECT name, token, step, error, status FROM jobs WHERE token = "c639b525-1ab1-44e6-bde3-96238cf13f2f"`)
	s.Require().Nil(err)

	s.Require().Len(jobs, 1)
	s.Require().ElementsMatch([]repo.Job{
		{
			Name:  "service1GetUser",
			Token: "c639b525-1ab1-44e6-bde3-96238cf13f2f",
			Step:  0,
			Error: sql.NullString{
				Valid:  true,
				String: "test error",
			},
			Status: int(repo.JobStatusFailed),
		},
	}, jobs)
}

func (s *Suite) Test_if_a_successfull_step_is_not_the_requests_last_step_increment_its_step_and_insert_the_new_steps_jobs() {
	s.importFile("delete_with_completed_jobs.sql")

	err := processInProgressRequests(&application.App{
		Repo: s.repo,
	})
	s.Require().Nil(err)

	var requests []repo.Request
	err = s.db.Select(&requests, "SELECT token, action, step, status FROM requests WHERE status = 1")
	s.Require().Nil(err)

	s.Require().Len(requests, 1)
	s.Require().ElementsMatch([]repo.Request{
		{
			Token:  "039a4e90-107b-4d7f-97f7-e1ad84316119",
			Action: string(types.ActionDelete),
			Step:   1,
			Status: int(repo.RequestStatusInProgress),
		},
	}, requests)

	var jobs []repo.Job
	err = s.db.Select(&jobs, `SELECT name, token, step, error, status FROM jobs WHERE token = "039a4e90-107b-4d7f-97f7-e1ad84316119"`)
	s.Require().Nil(err)

	s.Require().Len(jobs, 4)

	s.Require().ElementsMatch([]repo.Job{
		{
			Name:   "service1GetUser",
			Token:  "039a4e90-107b-4d7f-97f7-e1ad84316119",
			Step:   0,
			Error:  sql.NullString{},
			Status: int(repo.JobStatusCompleted),
		},
		{
			Name:   "service1DeleteUser",
			Token:  "039a4e90-107b-4d7f-97f7-e1ad84316119",
			Step:   1,
			Error:  sql.NullString{},
			Status: int(repo.JobStatusNew),
		},
		{
			Name:   "service2DeleteUser",
			Token:  "039a4e90-107b-4d7f-97f7-e1ad84316119",
			Step:   1,
			Error:  sql.NullString{},
			Status: int(repo.JobStatusNew),
		},
		{
			Name:   "service3DeleteUser",
			Token:  "039a4e90-107b-4d7f-97f7-e1ad84316119",
			Step:   1,
			Error:  sql.NullString{},
			Status: int(repo.JobStatusNew),
		},
	}, jobs)
}

func (s *Suite) Test_if_a_successfull_step_is_the_requests_last_step_update_the_requests_status_as_completed() {
	s.importFile("delete_with_completed_jobs_last_step.sql")

	err := processInProgressRequests(&application.App{
		Repo: s.repo,
	})
	s.Require().Nil(err)

	var requests []repo.Request
	err = s.db.Select(&requests, "SELECT token, action, step, status FROM requests WHERE status = 2")
	s.Require().Nil(err)

	s.Require().Len(requests, 1)
	s.Require().ElementsMatch([]repo.Request{
		{
			Token:  "80498d81-8de4-41fb-b1a5-53180cd56d73",
			Action: string(types.ActionDelete),
			Step:   1,
			Status: int(repo.RequestStatusCompleted),
		},
	}, requests)

	var jobs []repo.Job
	err = s.db.Select(&jobs, `SELECT name, token, step, error, status FROM jobs WHERE token = "80498d81-8de4-41fb-b1a5-53180cd56d73"`)
	s.Require().Nil(err)

	s.Require().Len(jobs, 4)

	s.Require().ElementsMatch([]repo.Job{
		{
			Name:   "service1GetUser",
			Token:  "80498d81-8de4-41fb-b1a5-53180cd56d73",
			Step:   0,
			Error:  sql.NullString{},
			Status: int(repo.JobStatusCompleted),
		},
		{
			Name:   "service1DeleteUser",
			Token:  "80498d81-8de4-41fb-b1a5-53180cd56d73",
			Step:   1,
			Error:  sql.NullString{},
			Status: int(repo.JobStatusCompleted),
		},
		{
			Name:   "service2DeleteUser",
			Token:  "80498d81-8de4-41fb-b1a5-53180cd56d73",
			Step:   1,
			Error:  sql.NullString{},
			Status: int(repo.JobStatusCompleted),
		},
		{
			Name:   "service3DeleteUser",
			Token:  "80498d81-8de4-41fb-b1a5-53180cd56d73",
			Step:   1,
			Error:  sql.NullString{},
			Status: int(repo.JobStatusCompleted),
		},
	}, jobs)
}
