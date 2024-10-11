package jobs

import (
	"fmt"

	"github.com/TheFranMan/tasker-common/types"

	"worker/application"
)

func Service1GetUser(app *application.App, request types.Request) (types.Extras, error) {
	// Call external API
	user, err := app.Service1.UserGet(request.Params.ID)
	if nil != err {
		return nil, types.Failure{Err: fmt.Errorf("cannot save email in extra: %w", err)}
	}

	return types.Extras{
		"email": user.Email,
	}, nil
}
