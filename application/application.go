package application

import (
	"github.com/robfig/cron/v3"

	"taskWorker/common"
	"taskWorker/repo"
	"taskWorker/service1"
)

type App struct {
	Config   *common.Config
	Repo     repo.Interface
	Cron     *cron.Cron
	Service1 service1.Interface
}
