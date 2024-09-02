package common

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env"

	_ "github.com/joho/godotenv/autoload"
)

type Envs struct {
	isLocal bool
	isStage bool
	isProd  bool
}

type Config struct {
	Envs
	Port int    `env:"PORT"`
	Env  string `env:"ENV"`

	DbUser string `env:"DB_USER"`
	DbPass string `env:"DB_PASS"`
	DbHost string `env:"DB_HOST"`
	DbPort string `env:"DB_PORT"`
	DbName string `env:"DB_NAME"`
}

func GetConfig() (*Config, error) {
	var config Config
	err := env.Parse(&config)
	if nil != err {
		return nil, fmt.Errorf("cannot parse env variables: %w", err)
	}

	config.setEnv()

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
