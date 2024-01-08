package main

import (
	"eljur/internal/config"
	"eljur/internal/logger"
	"fmt"
)

const configPath = "config/config.yaml"

func main() {
	cnf, err := config.GetConfig(configPath)
	if err != nil {
		panic(err)
	}
	l, err := logger.SetupLogger(&cnf.Log)
	if err != nil {
		panic(err)
	}
	l.Info(fmt.Sprintf("cnf: %+v", cnf))
	
}
