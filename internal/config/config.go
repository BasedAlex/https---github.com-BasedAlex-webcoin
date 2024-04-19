package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Env *EnvSetting
}

type EnvSetting struct {
	DBConn   string `env:"DB_CONNECTION"`
	DBCancel int    `env:"DB_CANCEL"`
	Port     string `env:"PORT"`
}

func New() *Config {
	e := &EnvSetting{}

	err := cleanenv.ReadConfig(".env", e)
	if err != nil {
		logrus.Panicf("read env config failed: %s", err)
	}

	return &Config{
		Env: e,
	}
}
