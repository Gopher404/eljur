package logger

import (
	"eljur/internal/config"
	"fmt"
	"io"
	"log/slog"
	"os"
)

func SetupLogger(cnf *config.LogConfig) (*slog.Logger, error) {
	var out io.Writer

	if cnf.Out == "cmd" {
		out = os.Stdout
	} else {
		f, err := os.Open(cnf.Out)
		if err != nil {
			return nil, err
		}
		out = f
	}

	var opts slog.HandlerOptions
	opts.AddSource = true

	switch cnf.Level {
	case "debug":
		opts.Level = slog.LevelDebug
		break
	case "error":
		opts.Level = slog.LevelError
		break
	case "info":
		opts.Level = slog.LevelInfo
		break
	case "warn":
		opts.Level = slog.LevelWarn
		break
	default:
		return nil, fmt.Errorf("unknown level: %s", cnf.Level)
	}

	var handler slog.Handler

	switch cnf.Type {
	case "text":
		handler = slog.NewTextHandler(out, &opts)
		break
	case "json":
		handler = slog.NewJSONHandler(out, &opts)
		break
	default:
		return nil, fmt.Errorf("unknown type: %s", cnf.Type)
	}

	l := slog.New(handler)

	return l, nil
}
