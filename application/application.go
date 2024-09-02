package application

import (
	repo "taskWorker/Repo"
	"taskWorker/common"
)

type App struct {
	Config *common.Config
	Repo   repo.Interface
}
