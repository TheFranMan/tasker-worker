//go:build integration

package main

import (
	"database/sql"

	"github.com/TheFranMan/tasker-common/types"

	"worker/application"
	"worker/service1"
)

func (s *Suite) Test_jobs_can_update_their_request_extras_column() {
	s.importFile("delete_inprogress.sql")

	mockService1 := new(service1.Mock)
	mockService1.On("UserGet", 1).Return(&service1.User{
		Email: "example_1@example.com",
	}, nil)

	err := processNewJobs(&application.App{
		Repo:     s.repo,
		Service1: mockService1,
	})
	s.Require().Nil(err)

	var requests []types.Request
	err = s.db.Select(&requests, "SELECT token, action, extras, step, status FROM requests WHERE status = 1")
	s.Require().Nil(err)

	s.Require().Len(requests, 1)
	s.Require().ElementsMatch([]types.Request{
		{
			Token:  "482a2d88-d38a-4509-ac94-beadff53c053",
			Action: string(types.ActionDelete),
			Extras: sql.NullString{
				Valid:  true,
				String: `{"email": "example_1@example.com"}`,
			},
			Step:   0,
			Status: 1,
		},
	}, requests)

	mockService1.AssertExpectations(s.T())
}
