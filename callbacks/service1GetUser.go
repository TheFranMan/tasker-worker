package callbacks

import (
	"fmt"
	"taskWorker/application"
	"taskWorker/repo"
)

func Service1GetUser(app *application.App, request repo.Request, id int) error {
	fmt.Println("In Service1GetUser lala:")
	// Call external API
	user, err := app.Service1.UserGet(id)
	if nil != err {
		return err
	}
	fmt.Printf("%+v\n", user)
	email := "example@example.com"

	// Save email to extras
	err = app.Repo.SaveExtra("email", email, request.Token)
	if nil != err {
		return fmt.Errorf("cannot save email in extra: %w", err)
	}

	// Mark job as completed
	return app.Repo.MarkJobCompleted(id)
}
