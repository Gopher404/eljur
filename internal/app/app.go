package app

import (
	"eljur/internal/config"
	"eljur/internal/delivery/httpServer"
	"eljur/internal/service/grades"
	"eljur/internal/service/subjects"
	"eljur/internal/service/users"
	"eljur/internal/storage"
	"log/slog"
)

func Run(cnf *config.Config, l *slog.Logger) error {
	s, err := storage.New(&cnf.DB)
	if err != nil {
		return err
	}
	l.Info("setup storage")

	gradesService := grades.New(s.Grades, s.Subjects)
	usersService := users.New(s.Users)
	subjectsService := subjects.New(s.Subjects)
	l.Info("setup services")

	handler := httpServer.NewHandler(l, gradesService, usersService, subjectsService)
	server := httpServer.NewServer(handler.GetMuxRouter(), &cnf.Bind)

	l.Info("run server")
	if err := server.Run(); err != nil {
		return err
	}
	return nil
}
