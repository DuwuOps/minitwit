package utils

import (
	"log/slog"
	"os"
)

func SetSlogDefaults() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    slog.SetDefault(logger)
}
