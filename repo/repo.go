package repo

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"taskWorker/common"
)

type Interface interface {
	GetAll() ([]Request, error)
}

type Request struct {
	Token        string         `db:"token"`
	RequestToken string         `db:"request_token"`
	Action       string         `db:"action"`
	Params       string         `db:"params"`
	Extras       sql.NullString `db:"extras"`
	Controller   string         `db:"controller"`
	Step         int            `db:"step"`
	Status       int            `db:"status"`
	Created      time.Time      `db:"created"`
	Completed    sql.NullTime   `db:"completed"`
}

type Repo struct {
	db *sqlx.DB
}

func NewRepo(config *common.Config) (*Repo, error) {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.DbUser, config.DbPass, config.DbHost, config.DbPort, config.DbName))
	if nil != err {
		return nil, err
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(5 * time.Minute)

	return &Repo{db}, nil
}

func (r *Repo) GetAll() ([]Request, error) {
	var requests []Request
	err := r.db.Select(&requests, "SELECT token, request_token, action, params, extras, controller, step, status, created, completed FROM requests")
	if nil != err {
		return nil, err
	}

	return requests, nil
}
