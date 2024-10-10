package common

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

type Envs struct {
	isLocal bool
	isStage bool
	isProd  bool
}

type Config struct {
	Envs
	Port int    `env:"PORT,notEmpty"`
	Env  string `env:"ENV"`

	DbUser string `env:"DB_USER,notEmpty"`
	DbPass string `env:"DB_PASS,notEmpty"`
	DbHost string `env:"DB_HOST,notEmpty"`
	DbPort string `env:"DB_PORT,notEmpty"`
	DbName string `env:"DB_NAME,notEmpty"`

	WrkEnabled               bool   `env:"WRK_ENABLED" envDefault:"true"`
	WrkJobCron               string `env:"WRK_JOB_CRON" envDefault:"*/15 * * * *"`
	WrkRequestNewCron        string `env:"WRK_REQUEST_NEW_CRON" envDefault:"*/5 * * * *"`
	WrkRequestInProgressCron string `env:"WRK_REQUEST_INPROGRESS_CRON" envDefault:"*/1 * * * *"`

	URLService1 url.URL `env:"URL_SERVICE1,notEmpty"`
	URLService2 url.URL `env:"URL_SERVICE2,notEmpty"`
	URLService3 url.URL `env:"URL_SERVICE3,notEmpty"`
}

func GetConfig() (*Config, error) {
	var config Config
	err := env.Parse(&config)
	if nil != err {
		return nil, fmt.Errorf("cannot parse env variables: %w", err)
	}

	config.setEnv()
	config.setLogging()

	return &config, nil
}

func (c *Config) setEnv() {
	if strings.HasPrefix(strings.ToLower(c.Env), "prod") {
		c.isLocal = false
		c.isStage = false
		c.isProd = true
		return
	}

	if strings.HasPrefix(strings.ToLower(c.Env), "stag") {
		c.isLocal = false
		c.isStage = true
		c.isProd = false
		return
	}

	c.isLocal = true
	c.isStage = false
	c.isProd = false
}

func (c *Config) setLogging() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	log.SetLevel(log.InfoLevel)

	if c.isLocal {
		log.SetLevel(log.DebugLevel)
	}
}
