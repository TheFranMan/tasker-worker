package jobs

import (
	"fmt"

	"github.com/TheFranMan/tasker-common/types"
	log "github.com/sirupsen/logrus"

	"worker/application"
)

func Service2UpdateUser(app *application.App, request types.Request) (types.Extras, error) {
	err := app.Service2.UpdateUser(request.Params.ID, request.Params.Email)
	if nil != err {
		log.WithField("id", request.Params.ID).WithError(err).Error("cannot update user email from service 2")
		return nil, types.Failure{Err: fmt.Errorf("cannot update user email from service 2: %w", err)}
	}

	return nil, nil
}
