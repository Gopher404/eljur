package suite

import (
	"eljur/internal/config"
	"eljur/internal/storage"
	"fmt"
)

const configPath = "../config/config.yaml"

func GetStorage() (*storage.Storage, error) {
	cnf, err := config.GetConfig(configPath)
	if err != nil {
		return nil, err
	}

	s, err := storage.New("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", cnf.DB.User, cnf.DB.Password, cnf.DB.Host, cnf.DB.Schema))
	if err != nil {
		return nil, err
	}

	return s, nil
}
