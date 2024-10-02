package main

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	db       *sqlx.DB
	pool     *dockertest.Pool
	resource *dockertest.Resource
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

	err = s.pool.Retry(func() error {
		var err error
		s.db, err = sqlx.Open("mysql", fmt.Sprintf("root:secret@(localhost:%s)/mysql", s.resource.GetPort("3306/tcp")))
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
}

func (s *Suite) TearDownSuite() {
	err := s.pool.Purge(s.resource)
	if nil != err {
		s.FailNowf(err.Error(), "cannot purge dockertest mysql resource")
	}
}

func (s *Suite) AfterTest() {
	s.T().Log("In Aftertest")
}

func (s *Suite) Test_new_requests_are_marked_in_progress_and_their_first_step_jobs_inserted() {
	s.T().Log("In test1")

}

func (s *Suite) Test_if_the_request_is_new__insert_the_first_step_jobs_and_update_the_request_status_as_in_progress() {
	s.T().Log("In test2")

}

func (s *Suite) Test_if_the_request_has_an_errored_job_reinsert_the_job_and_do_not_update_the_request_status() {
	s.T().Log("In test3")

}

func (s *Suite) Test_if_a_successfull_step_is_not_the_requests_last_step_increment_its_step_and_insert_the_new_steps_jobs() {
	s.T().Log("In test4")

}

func (s *Suite) Test_if_a_successfull_step_is_the_requests_last_step_update_the_requests_status_as_completed() {
	s.T().Log("In test5")

}
