package application

import (
	"github.com/robfig/cron/v3"

	"worker/common"
	"worker/repo"
	"worker/service1"
)

type App struct {
	Config   *common.Config
	Repo     repo.Interface
	Cron     *cron.Cron
	Service1 service1.Interface
}
