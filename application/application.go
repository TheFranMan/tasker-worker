package application

import (
	"github.com/robfig/cron/v3"

	"worker/common"
	"worker/repo"
	"worker/service1"
	"worker/service2"
)

type App struct {
	Config   *common.Config
	Repo     repo.Interface
	Cron     *cron.Cron
	Service1 service1.Interface
	Service2 service2.Interface
}
