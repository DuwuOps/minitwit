# Coding Conventions

## Logging

When logging or printing, use the "log/slog"-package and not "log" or "fmt".

### When to use what

If there is a context, please use a context based method, such as `slog.InfoContext` or `slog.ErrorContext`. Otherwise, do not feel the need to create a context. Just use non-context based methods such as `slog.Info` or `slog.Error`.

Logging helper-methods go in `src/utils/logger.go`. Here already exists some helper-methods for echo-contexts, such as `InfoEchoContext`, such that there is no need for extracting the`context.Context`-object from the `echo.Context` everytime.