package scimpatch

import "context"

type PatcherLogger interface {
	Error(args ...interface{})
	Debug(args ...interface{})
}

type noopPatcherLogger struct{}

func (l noopPatcherLogger) Error(args ...interface{}) {}
func (l noopPatcherLogger) Debug(args ...interface{}) {}

var noop = noopPatcherLogger{}

type loggerKey struct{}

// key is the context key for the logger.
var key = loggerKey{}

func AddLogger(ctx context.Context, logger *PatcherLogger) context.Context {
	return context.WithValue(ctx, key, logger)
}

func logger(ctx context.Context) PatcherLogger {
	l, ok := ctx.Value(key).(PatcherLogger)
	if !ok {
		return noop
	}
	return l
}
