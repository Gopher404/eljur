package main

import (
	"eljur/internal/app"
	"eljur/internal/config"
	"eljur/internal/pkg/logger"
)

const configPath = "config/config.yaml"

func main() {
	cnf, err := config.GetConfig(configPath)
	if err != nil {
		panic(err)
	}
	l, err := logger.SetupLogger(cnf.Log)
	if err != nil {
		panic(err)
	}

	if err := app.Run(cnf, l); err != nil {
		panic(err)
	}
}
