package callbacks

import (
	"fmt"
	"worker/application"
	"worker/repo"

	log "github.com/sirupsen/logrus"
)

func Service2DeleteUser(app *application.App, request repo.Request, id int) error {
	err := app.Service2.DeleteUser(request.Params.ID)
	if nil != err {
		log.WithField("id", request.Params.ID).WithError(err).Error("cannot delete user from service 2")
		return app.Repo.MarkJobFailed(id, fmt.Errorf("cannot delete user from service 2: %w", err))
	}

	return app.Repo.MarkJobCompleted(id)
}
