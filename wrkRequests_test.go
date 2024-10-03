package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

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

	s.resource.Expire(60 * 1)
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
			Token:  "test-token-1",
			Action: "Delete",
			Step:   0,
			Status: 1,
		},
		{
			Token:  "test-token-2",
			Action: "Delete",
			Step:   0,
			Status: 1,
		},
	}, requests)

	var jobs []repo.Job
	err = s.db.Select(&jobs, "SELECT name, token, step, error, status FROM jobs")
	s.Require().Nil(err)

	s.Require().Len(jobs, 2)
	s.Require().ElementsMatch([]repo.Job{
		{
			Name:   "service1GetUser",
			Token:  "test-token-1",
			Step:   0,
			Error:  sql.NullString{},
			Status: 0,
		},
		{
			Name:   "service1GetUser",
			Token:  "test-token-2",
			Step:   0,
			Error:  sql.NullString{},
			Status: 0,
		},
	}, jobs)
}

func (s *Suite) Test_if_the_request_has_an_errored_job_reinsert_the_job_and_do_not_update_the_request_status() {
	s.importFile("delete_with_error_jobs.sql")

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
			Token:  "test-token-1",
			Action: "Delete",
			Step:   0,
			Status: 1,
		},
		{
			Token:  "test-token-2",
			Action: "Delete",
			Step:   0,
			Status: 1,
		},
	}, requests)

	var jobs []repo.Job
	err = s.db.Select(&jobs, "SELECT name, token, step, error, status FROM jobs")
	s.Require().Nil(err)

	s.Require().Len(jobs, 4)
	s.Require().ElementsMatch([]repo.Job{
		{
			Name:  "service1GetUser",
			Token: "test-token-1",
			Step:  0,
			Error: sql.NullString{
				Valid:  true,
				String: "test error",
			},
			Status: 4,
		},
		{
			Name:  "service1GetUser",
			Token: "test-token-2",
			Step:  0,
			Error: sql.NullString{
				Valid:  true,
				String: "test error",
			},
			Status: 4,
		},
		{
			Name:   "service1GetUser",
			Token:  "test-token-1",
			Step:   0,
			Error:  sql.NullString{},
			Status: 0,
		},
		{
			Name:   "service1GetUser",
			Token:  "test-token-2",
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

	s.Require().Len(requests, 2)
	s.Require().ElementsMatch([]repo.Request{
		{
			Token:  "test-token-1",
			Action: "Delete",
			Step:   0,
			Status: 3,
		},
		{
			Token:  "test-token-2",
			Action: "Delete",
			Step:   0,
			Status: 3,
		},
	}, requests)

	var jobs []repo.Job
	err = s.db.Select(&jobs, "SELECT name, token, step, error, status FROM jobs")
	s.Require().Nil(err)

	s.Require().Len(jobs, 2)
	s.Require().ElementsMatch([]repo.Job{
		{
			Name:  "service1GetUser",
			Token: "test-token-1",
			Step:  0,
			Error: sql.NullString{
				Valid:  true,
				String: "test error",
			},
			Status: 3,
		},
		{
			Name:  "service1GetUser",
			Token: "test-token-2",
			Step:  0,
			Error: sql.NullString{
				Valid:  true,
				String: "test error",
			},
			Status: 3,
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

	s.Require().Len(requests, 2)
	s.Require().ElementsMatch([]repo.Request{
		{
			Token:  "test-token-1",
			Action: "Delete",
			Step:   1,
			Status: 1,
		},
		{
			Token:  "test-token-2",
			Action: "Delete",
			Step:   1,
			Status: 1,
		},
	}, requests)

	var jobs []repo.Job
	err = s.db.Select(&jobs, "SELECT name, token, step, error, status FROM jobs")
	s.Require().Nil(err)

	s.Require().Len(jobs, 8)

	s.Require().ElementsMatch([]repo.Job{
		{
			Name:   "service1GetUser",
			Token:  "test-token-1",
			Step:   0,
			Error:  sql.NullString{},
			Status: 2,
		},
		{
			Name:   "service1GetUser",
			Token:  "test-token-2",
			Step:   0,
			Error:  sql.NullString{},
			Status: 2,
		},
		{
			Name:   "service1DeleteUser",
			Token:  "test-token-1",
			Step:   1,
			Error:  sql.NullString{},
			Status: 0,
		},
		{
			Name:   "service2DeleteUser",
			Token:  "test-token-1",
			Step:   1,
			Error:  sql.NullString{},
			Status: 0,
		},
		{
			Name:   "service3DeleteUser",
			Token:  "test-token-1",
			Step:   1,
			Error:  sql.NullString{},
			Status: 0,
		},
		{
			Name:   "service1DeleteUser",
			Token:  "test-token-2",
			Step:   1,
			Error:  sql.NullString{},
			Status: 0,
		},
		{
			Name:   "service2DeleteUser",
			Token:  "test-token-2",
			Step:   1,
			Error:  sql.NullString{},
			Status: 0,
		},
		{
			Name:   "service3DeleteUser",
			Token:  "test-token-2",
			Step:   1,
			Error:  sql.NullString{},
			Status: 0,
		},
	}, jobs)
}

func (s *Suite) Test_if_a_successfull_step_is_the_requests_last_step_update_the_requests_status_as_completed() {
	s.importFile("update_email_with_completed_jobs.sql")

	err := processInProgressRequests(&application.App{
		Repo: s.repo,
	})
	s.Require().Nil(err)

	var requests []repo.Request
	err = s.db.Select(&requests, "SELECT token, action, step, status FROM requests WHERE status = 2")
	s.Require().Nil(err)

	s.Require().Len(requests, 2)
	s.Require().ElementsMatch([]repo.Request{
		{
			Token:  "test-token-1",
			Action: "update_email",
			Step:   0,
			Status: 2,
		},
		{
			Token:  "test-token-2",
			Action: "update_email",
			Step:   0,
			Status: 2,
		},
	}, requests)

	var jobs []repo.Job
	err = s.db.Select(&jobs, "SELECT name, token, step, error, status FROM jobs")
	s.Require().Nil(err)

	s.Require().Len(jobs, 6)

	s.Require().ElementsMatch([]repo.Job{
		{
			Name:   "service1UpdateUser",
			Token:  "test-token-1",
			Step:   0,
			Error:  sql.NullString{},
			Status: 2,
		},
		{
			Name:   "service2UpdateUser",
			Token:  "test-token-1",
			Step:   0,
			Error:  sql.NullString{},
			Status: 2,
		},
		{
			Name:   "service3UpdateUser",
			Token:  "test-token-1",
			Step:   0,
			Error:  sql.NullString{},
			Status: 2,
		},
		{
			Name:   "service1UpdateUser",
			Token:  "test-token-2",
			Step:   0,
			Error:  sql.NullString{},
			Status: 2,
		},
		{
			Name:   "service2UpdateUser",
			Token:  "test-token-2",
			Step:   0,
			Error:  sql.NullString{},
			Status: 2,
		},
		{
			Name:   "service3UpdateUser",
			Token:  "test-token-2",
			Step:   0,
			Error:  sql.NullString{},
			Status: 2,
		},
	}, jobs)
}
