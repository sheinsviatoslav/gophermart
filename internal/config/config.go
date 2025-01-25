package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

var (
	RunAddress           = flag.String("a", ":8000", "server address")
	DatabaseURI          = flag.String("d", "database=postgres", "database URI")
	AccrualSystemAddress = flag.String("r", "http://localhost:8080", "accrual system address")
)

func Init() {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	if cfg.RunAddress != "" {
		*RunAddress = cfg.RunAddress
	}

	if cfg.DatabaseURI != "" {
		*DatabaseURI = cfg.DatabaseURI
	}

	if cfg.AccrualSystemAddress != "" {
		*AccrualSystemAddress = cfg.AccrualSystemAddress
	}
}
