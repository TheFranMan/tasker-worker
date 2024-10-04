package callbacks

import (
	"fmt"

	"worker/application"
	"worker/repo"

	"github.com/TheFranMan/tasker-common/types"
)

func Service1GetUser(app *application.App, request repo.Request) error {
	// Call external API
	user, err := app.Service1.UserGet(request.Params.ID)
	if nil != err {
		return types.Failure{Err: fmt.Errorf("cannot save email in extra: %w", err)}
	}

	// Save email to extras
	err = app.Repo.SaveExtra("email", user.Email, request.Token)
	if nil != err {
		return types.Failure{Err: fmt.Errorf("cannot save email in extra: %w", err)}
	}

	return nil
}
