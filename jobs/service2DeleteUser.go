package jobs

import (
	"fmt"

	"github.com/TheFranMan/tasker-common/types"
	log "github.com/sirupsen/logrus"

	"worker/application"
)

func Service2DeleteUser(app *application.App, request types.Request) (types.Extras, error) {
	err := app.Service2.DeleteUser(request.Params.ID)
	if nil != err {
		log.WithField("id", request.Params.ID).WithError(err).Error("cannot delete user from service 2")
		return nil, types.Failure{Err: fmt.Errorf("cannot delete user from service 2: %w", err)}
	}

	return nil, nil
}
