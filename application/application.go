package application

import (
	"github.com/robfig/cron/v3"

	"taskWorker/common"
	"taskWorker/repo"
)

type App struct {
	Config *common.Config
	Repo   repo.Interface
	Cron   *cron.Cron
}
