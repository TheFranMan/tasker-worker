package callbacks

import (
	"encoding/json"
	"fmt"
	"worker/application"
	"worker/repo"

	log "github.com/sirupsen/logrus"
)

func Service3DeleteUser(app *application.App, request repo.Request, id int) error {
	var extras map[string]any
	err := json.Unmarshal([]byte(request.Extras.String), &extras)
	if nil != err {
		log.WithError(err).Error("cannot unmarshall extras")
		return app.Repo.MarkJobFailed(id, fmt.Errorf("cannot unmarshal extras from service 3: %w", err))
	}

	err = app.Service3.DeleteUser(extras["email"].(string))
	if nil != err {
		log.WithError(err).Error("cannot delete user from service 3")
		return app.Repo.MarkJobFailed(id, fmt.Errorf("cannot delete user from service 3: %w", err))
	}

	return app.Repo.MarkJobCompleted(id)
}
