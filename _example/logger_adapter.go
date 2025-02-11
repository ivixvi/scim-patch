package main

import "log"

type PatcherLogger struct {
	logger *log.Logger
}

func newLogger(logger *log.Logger) PatcherLogger {
	return PatcherLogger{
		logger: logger,
	}
}

func (l PatcherLogger) Error(args ...interface{}) {
	l.logger.Println(args...)
}

func (l PatcherLogger) Debug(args ...interface{}) {
	l.logger.Println(args...)
}
