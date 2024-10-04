package callbacks

import (
	"encoding/json"
	"fmt"

	"github.com/TheFranMan/tasker-common/types"
	log "github.com/sirupsen/logrus"

	"worker/application"
	"worker/repo"
)

func Service3DeleteUser(app *application.App, request repo.Request) error {
	var extras map[string]any
	err := json.Unmarshal([]byte(request.Extras.String), &extras)
	if nil != err {
		log.WithError(err).Error("cannot unmarshall extras")
		return types.Failure{Err: fmt.Errorf("cannot unmarshal extras from service 3: %w", err)}
	}

	err = app.Service3.DeleteUser(extras["email"].(string))
	if nil != err {
		log.WithError(err).Error("cannot delete user from service 3")
		return types.Failure{Err: fmt.Errorf("cannot delete user from service 3: %w", err)}
	}

	return nil
}
