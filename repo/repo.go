package repo

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"taskWorker/common"
)

type Interface interface{}

type Repo struct {
	db *sqlx.DB
}

func NewRepo(config *common.Config) (*Repo, error) {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DbUser, config.DbPass, config.DbHost, config.DbPort, config.DbName))
	if nil != err {
		return nil, err
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(5 * time.Minute)

	return &Repo{db}, nil
}
