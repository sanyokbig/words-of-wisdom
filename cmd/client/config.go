package main

import "github.com/caarlos0/env/v6"

type Config struct {
	ServerAddr string `env:"SERVER_ADDR"`
}

func parse() (*Config, error) {
	config := &Config{}

	err := env.Parse(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
