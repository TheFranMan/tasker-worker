package callbacks

import (
	"fmt"
	"worker/application"
	"worker/repo"
)

func Service1GetUser(app *application.App, request repo.Request, id int) error {
	// Call external API
	user, err := app.Service1.UserGet(id)
	if nil != err {
		return err
	}

	// Save email to extras
	err = app.Repo.SaveExtra("email", user.Email, request.Token)
	if nil != err {
		return fmt.Errorf("cannot save email in extra: %w", err)
	}

	// Mark job as completed
	return app.Repo.MarkJobCompleted(id)
}
