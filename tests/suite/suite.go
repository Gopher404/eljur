package suite

import (
	"eljur/internal/config"
	"eljur/internal/service/grades"
	"eljur/internal/storage"
	"fmt"
	"log/slog"
	"os"
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

func GetGradesService() (*grades.GradeService, error) {
	s, err := GetStorage()
	if err != nil {
		return nil, err
	}
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return grades.New(s.Grades, s.Subjects, l), nil

}
