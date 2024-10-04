package main

import (
	"database/sql"
	"worker/application"
	"worker/repo"
	"worker/service1"
)

func (s *Suite) Test_jobs_can_update_their_request_extras_column() {
	s.importFile("delete_email_inprogress.sql")

	mockService1 := new(service1.Mock)
	mockService1.On("UserGet", 1).Return(&service1.User{
		Email: "example_1@example.com",
	}, nil)
	mockService1.On("UserGet", 2).Return(&service1.User{
		Email: "example_2@example.com",
	}, nil)

	err := processNewJobs(&application.App{
		Repo:     s.repo,
		Service1: mockService1,
	})
	s.Require().Nil(err)

	var requests []repo.Request
	err = s.db.Select(&requests, "SELECT token, action, extras, step, status FROM requests WHERE status = 1")
	s.Require().Nil(err)

	s.Require().Len(requests, 2)
	s.Require().ElementsMatch([]repo.Request{
		{
			Token:  "test-token-1",
			Action: "Delete",
			Extras: sql.NullString{
				Valid:  true,
				String: `{"email": "example_1@example.com"}`,
			},
			Step:   0,
			Status: 1,
		},
		{
			Token:  "test-token-2",
			Action: "Delete",
			Extras: sql.NullString{
				Valid:  true,
				String: `{"email": "example_2@example.com"}`,
			},
			Step:   0,
			Status: 1,
		},
	}, requests)

	mockService1.AssertExpectations(s.T())
}
