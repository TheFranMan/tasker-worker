package callbacks

import (
	"fmt"
	"worker/application"
	"worker/repo"

	log "github.com/sirupsen/logrus"
)

func Service1DeleteUser(app *application.App, request repo.Request, id int) error {
	err := app.Service1.DeleteUser(request.Params.ID)
	if nil != err {
		log.WithField("id", request.Params.ID).WithError(err).Error("cannot delete user from service 1")
		return app.Repo.MarkJobFailed(id, fmt.Errorf("cannot delete user from service 1: %w", err))
	}

	return app.Repo.MarkJobCompleted(id)
}
