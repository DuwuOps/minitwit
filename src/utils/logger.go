package utils

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
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


func LogErrorEchoContext(echoCtx echo.Context, msg string, err error) {
	ctx:= echoCtx.Request().Context()
	LogErrorContext(ctx, msg, err)
}


func LogRouteStart(echoCtx echo.Context, routeName string, route string) {
	ctx:= echoCtx.Request().Context()
	logMsg := fmt.Sprintf("🎺 User entered %s via route \"%s\"", routeName, route)
	slog.ErrorContext(ctx, logMsg,  slog.Any("HTTP method", echoCtx.Request().Method))
}


func InfoEchoContext(echoCtx echo.Context, logMsg string) {
	ctx:= echoCtx.Request().Context()
	slog.InfoContext(ctx, logMsg)
}