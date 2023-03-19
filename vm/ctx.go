package vm

import (
	"context"
	"log"
)

type ctxKey uint

const (
	environmentKey ctxKey = iota
	loggerKey      ctxKey = iota
)

// EnvironmentFromContext returns the console environment from the context.
func EnvironmentFromContext(ctx context.Context) *Environment {
	return ctx.Value(environmentKey).(*Environment)
}

// LoggerFromContext returns the logger from the context.
func LoggerFromContext(ctx context.Context) *log.Logger {
	return ctx.Value(loggerKey).(*log.Logger)
}
