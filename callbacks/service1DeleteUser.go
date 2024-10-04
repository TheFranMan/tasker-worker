package callbacks

import (
	"fmt"

	"github.com/TheFranMan/tasker-common/types"
	log "github.com/sirupsen/logrus"

	"worker/application"
	"worker/repo"
)

func Service1DeleteUser(app *application.App, request repo.Request) error {
	err := app.Service1.DeleteUser(request.Params.ID)
	if nil != err {
		log.WithField("id", request.Params.ID).WithError(err).Error("cannot delete user from service 1")
		return types.Failure{Err: fmt.Errorf("cannot save email in extra: %w", err)}
	}

	return nil
}
