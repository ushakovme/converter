package infrastructure

import (
	"github.com/caarlos0/env"
	"github.com/pkg/errors"
	"regexp"
)

type Config struct {
	Port string `env:"CONVERTER_PORT" envDefault:":8000"`
}

func LoadConfig() (Config, error) {
	cnf := Config{}
	err := env.Parse(&cnf)

	if err != nil {
		return cnf, err
	}

	err = validateConfig(cnf)

	return cnf, err
}

func validateConfig(cnf Config) error {
	matched, err := regexp.MatchString(`:\d+`, cnf.Port)
	if err != nil {
		return errors.WithStack(err)
	}
	if !matched {
		return errors.New("port is not valid")
	}

	return nil
}
