package logging

import (
	"context"
)

type LoggingEnv int
type LogLevel int

const (
	// Development environment.
	development LoggingEnv = iota
	// Production environment.
	production

	// DebugLevel represents the debug logging level.
	debugLevel LogLevel = iota
	// InfoLevel represents the info logging level.
	infoLevel
	// WarnLevel represents the warn logging level.
	warnLevel
	// ErrorLevel represents the error logging level.
	errorLevel
	// FatalLevel represents the fatal logging level.
	fatalLevel
	// PanicLevel represents the panic logging level.
	panicLevel
)

// global var which define the global logging environment
var logEnv LoggingEnv = production
var serviceName string = "serviceName undefined"

// LoggingFacade provide access to the logger
type LoggingFacade struct {
	NewLogHandler func(LogLevel, ...context.Context) ILogHandler
	*logLevelHandler
}

func NewLoggingFacade(svcName ...string) *LoggingFacade {
	// LogEnv = Production
	if len(svcName) > 0 {
		serviceName = svcName[0]
	}
	return &LoggingFacade{
		NewLogHandler:   newLogHandler,
		logLevelHandler: NewLogLevelHandler(),
	}
}

func (lh *LoggingFacade) SetLoggingEnvToDevelopment() {
	logEnv = development
}

func (lh *LoggingFacade) SetLoggingEnvToProduction() {
	logEnv = production
}

func (lh *LoggingFacade) GetLoggingEnv() string {
	if logEnv == production {
		return "production"
	} else {
		return "development"
	}
}

func (lh *LoggingFacade) SetSvcName(name string) {
	serviceName = name
}

// LogLevelHandler handle the LogLevel
type logLevelHandler struct{}

func NewLogLevelHandler() *logLevelHandler {
	return &logLevelHandler{}
}

func (l *logLevelHandler) LLHDebug() LogLevel {
	return debugLevel
}
func (l *logLevelHandler) LLHInfo() LogLevel {
	return infoLevel
}
func (l *logLevelHandler) LLHWarn() LogLevel {
	return warnLevel
}
func (l *logLevelHandler) LLHError() LogLevel {
	return errorLevel
}
func (l *logLevelHandler) LLHFatal() LogLevel {
	return fatalLevel
}
func (l *logLevelHandler) LLHPanic() LogLevel {
	return panicLevel
}
