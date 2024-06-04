package app

import (
	"eljur/internal/config"
	"eljur/internal/delivery/httpHandler"
	"eljur/internal/delivery/httpServer"
	"eljur/internal/pkg/metric"
	"eljur/internal/service/grades"
	"eljur/internal/service/schedule"
	"eljur/internal/service/subjects"
	"eljur/internal/service/users"
	"eljur/internal/storage"
	"eljur/pkg/AuthClient"
	"log/slog"
)

func Run(cnf *config.Config, l *slog.Logger) error {
	s, err := storage.New(&cnf.DB)
	if err != nil {
		return err
	}
	l.Info("setup storage")

	authClient, err := AuthClient.New(cnf.SSO.Host, cnf.SSO.Port, cnf.SSO.AppKey)
	if err != nil {
		return err
	}

	gradesService := grades.New(s)
	subjectsService := subjects.New(s, gradesService)
	usersService := users.New(s, authClient, gradesService)
	scheduleService := schedules.New(s, s.Users, &cnf.Schedule)

	l.Info("setup services")

	handler := httpHandler.NewHandler(l,
		gradesService,
		usersService,
		subjectsService,
		scheduleService,
		authClient,
		cnf.Bind.HttpTimeOut,
	)

	server := httpServer.NewServer(handler.GetMuxRouter(), &cnf.Bind)

	metric.CountRPS()

	l.Info("run server")
	if err := server.Run(); err != nil {
		return err
	}
	return nil
}
