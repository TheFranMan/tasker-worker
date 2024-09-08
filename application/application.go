package application

import (
	"taskWorker/common"
	"taskWorker/repo"
)

type App struct {
	Config *common.Config
	Repo   repo.Interface
}
