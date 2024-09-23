package callbacks

import (
	"taskWorker/application"
	"taskWorker/repo"

	log "github.com/sirupsen/logrus"
)

func Service1DeleteUser(app *application.App, request repo.Request, id int) error {
	err := app.Service1.DeleteUser(id)
	if nil != err {
		log.WithError(err).Error("cannot delete user from service 1")
		return app.Repo.MarkJobError(id)
	}

	return app.Repo.MarkJobCompleted(id)
}
