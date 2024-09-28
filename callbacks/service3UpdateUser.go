package callbacks

import (
	"fmt"
	"worker/application"
	"worker/repo"

	log "github.com/sirupsen/logrus"
)

func Service3UpdateUser(app *application.App, request repo.Request, id int) error {
	err := app.Service3.UpdateUser(request.Params.Email, request.Params.Email)
	if nil != err {
		log.WithField("id", request.Params.ID).WithError(err).Error("cannot update user email from service 3")
		return app.Repo.MarkJobFailed(id, fmt.Errorf("cannot update user email from service 3: %w", err))
	}

	return app.Repo.MarkJobCompleted(id)
}
