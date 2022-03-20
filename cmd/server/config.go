package main

import "github.com/caarlos0/env/v6"

type Config struct {
	TCPPort        int    `env:"TCP_PORT"`
	QuotesFilePath string `env:"QUOTES_FILE_PATH"`
}

func parse() (*Config, error) {
	config := &Config{}

	err := env.Parse(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
