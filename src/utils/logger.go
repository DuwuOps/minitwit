package utils

import (
	"log/slog"
	"os"
)

func SetSlogDefaults() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    slog.SetDefault(logger)
}


func LogError(msg string, err error) {
	slog.Error(msg, slog.Any("error", err))
}


func LogFatal(msg string, err error) {
	LogError(msg, err)
	os.Exit(1)
}


func LogErrorContext(ctx context.Context, msg string, err error) {
	slog.ErrorContext(ctx, msg, slog.Any("error", err))
}

func InfoEchoContext(echoCtx echo.Context, logMsg string) {
	ctx:= echoCtx.Request().Context()
	slog.InfoContext(ctx, logMsg)
}