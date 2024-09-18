package callbacks

import (
	"fmt"
	"taskWorker/application"
	"taskWorker/repo"
)

func Service1GetUser(app *application.App, request repo.Request, id int) error {
	fmt.Println("In Service1GetUser")
	// Call external API
	email := "example@example.com"

	// Save email to extras
	err := app.Repo.SaveExtra("email", email, request.Token)
	if nil != err {
		return fmt.Errorf("cannot save email in extra: %w", err)
	}

	// Mark job as completed
	return app.Repo.MarkJobCompleted(id)
}
