package logging

import (
	"context"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"os"
	"time"
)

// depending of the environment, select if the log should be print or not
func isLoggingON(logLvl LogLevel) bool {
	if logEnv == production {
		if logLvl >= errorLevel {
			return true
		}
		return false
	} else {
		return true
	}
}

type ILogHandler interface {
	Str(k, v string) *LogHandler
	Int(k string, v int) *LogHandler
	Float64(k string, v float64) *LogHandler
	Err(error) *LogHandler
	Msg(msg string)
	Msgf(msg string, v interface{})
	Send()
}

var _ *LogHandler = new(LogHandler)

// LogHandler handle the logger
type LogHandler struct {
	// logger zerolog.Logger
	event *zerolog.Event
	llvl  LogLevel
}

func newLogHandler(llv LogLevel, ctx ...context.Context) ILogHandler {
	var logger *Logger
	if len(ctx) > 0 {
		logger = NewLoggerTrace(ctx[0])
	} else {
		logger = NewLogger()
	}

	var evt *zerolog.Event
	switch llv {
	case debugLevel:
		evt = logger.Debug()
	case infoLevel:
		evt = logger.Info()
	case warnLevel:
		evt = logger.Warn()
	case errorLevel:
		evt = logger.Error()
	case fatalLevel:
		evt = logger.Fatal()
	case panicLevel:
		evt = logger.Panic()
	}

	return &LogHandler{
		// logger: logger.log,
		event: evt,
		llvl:  llv}
}

func (l *LogHandler) Str(k, v string) *LogHandler {
	ok := isLoggingON(l.llvl)
	if ok {
		l.event.Str(k, v)
	}
	return l
}

func (l *LogHandler) Int(k string, v int) *LogHandler {
	ok := isLoggingON(l.llvl)
	if ok {
		l.event.Int(k, v)
	}
	return l
}

func (l *LogHandler) Float64(k string, v float64) *LogHandler {
	ok := isLoggingON(l.llvl)
	if ok {
		l.event.Float64(k, v)
	}
	return l
}

func (l *LogHandler) Err(err error) *LogHandler {
	ok := isLoggingON(l.llvl)
	if ok {
		l.event.Err(err)
	}
	return l
}

func (l *LogHandler) Msg(msg string) {
	ok := isLoggingON(l.llvl)
	if ok {
		l.event.Msg(msg)
	}
}

func (l *LogHandler) Msgf(msg string, v interface{}) {
	ok := isLoggingON(l.llvl)
	if ok {
		l.event.Msgf("%v: %v", msg, v)
	}
}

func (l *LogHandler) Send() {
	ok := isLoggingON(l.llvl)
	if ok {
		l.event.Send()
	}
}

// ILogger is an interface for the logger.
type ILogger interface {
	Debug() *zerolog.Event
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event
	Fatal() *zerolog.Event
	Panic() *zerolog.Event
}

type Logger struct {
	log zerolog.Logger
}

// NewNormalLogger creates a new zerolog Logger for normal logging.
func NewLogger() *Logger {
	if serviceName == "" {
		serviceName = "undefined serviceName"
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano
	l := zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Str("service", serviceName).
		Logger()

	return &Logger{l}
}

// Return a Logger set with a hook for the trace to records the log
func NewLoggerTrace(ctx context.Context) *Logger {
	if serviceName == "" {
		serviceName = "undefined serviceName"
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano
	l := zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Str("service", serviceName).
		Logger()

	return &Logger{l.Hook(zerologTraceHook(ctx))}
}

// Debug returns a Debug level event for normal logging.
func (l *Logger) Debug() *zerolog.Event {
	// l.log.Level(zerolog.DebugLevel)
	return l.log.Debug()
}

// Info returns an Info level event for normal logging.
func (l *Logger) Info() *zerolog.Event {
	// l.log.Level(zerolog.DebugInfo)
	return l.log.Info()
}

// Warn returns a Warn level event for normal logging.
func (l *Logger) Warn() *zerolog.Event {
	// l.log.Level(zerolog.DebugWarn)
	return l.log.Warn()
}

// Error returns an Error level event for normal logging.
func (l *Logger) Error() *zerolog.Event {
	// l.log.Level(zerolog.DebugError)
	return l.log.Error()
}

// Fatal returns a Fatal level event for normal logging.
func (l *Logger) Fatal() *zerolog.Event {
	// l.log.Level(zerolog.DebugFatal)
	return l.log.Fatal()
}

// Panic returns a Panic level event for normal logging.
func (l *Logger) Panic() *zerolog.Event {
	// l.log.Level(zerolog.DebugPanic)
	return l.log.Panic()
}

// zerologTraceHook is a hook that;
// (a) adds TraceIds & spanIds to logs of all LogLevels
// (b) adds logs to the active span as events.
func zerologTraceHook(ctx context.Context) zerolog.HookFunc {
	return func(e *zerolog.Event, level zerolog.Level, message string) {
		if level == zerolog.NoLevel {
			return
		}
		if !e.Enabled() {
			return
		}

		if ctx == nil {
			return
		}

		span := trace.SpanFromContext(ctx)
		if !span.IsRecording() {
			return
		}

		{ // (a) adds TraceIds & spanIds to logs.
			//
			// TODO: (komuw) add stackTraces maybe.
			//
			sCtx := span.SpanContext()
			if sCtx.HasTraceID() {
				e.Str("traceId", sCtx.TraceID().String())
			}
			if sCtx.HasSpanID() {
				e.Str("spanId", sCtx.SpanID().String())
			}
		}

		{ // (b) adds logs to the active span as events.
			if logEnv == production {
				// In production, only add logs for specified log levels
				if level >= zerolog.ErrorLevel {
					attrs := make([]attribute.KeyValue, 0)
					logSeverityKey := attribute.Key("log.severity")
					logMessageKey := attribute.Key("log.message")
					attrs = append(attrs, logSeverityKey.String(level.String()))
					attrs = append(attrs, logMessageKey.String(message))
					span.AddEvent("log", trace.WithAttributes(attrs...))
					if level >= zerolog.ErrorLevel {
						span.SetStatus(codes.Error, message)
					}
				}
			} else {
				// In development, add logs for all levels
				attrs := make([]attribute.KeyValue, 0)
				logSeverityKey := attribute.Key("log.severity")
				logMessageKey := attribute.Key("log.message")
				attrs = append(attrs, logSeverityKey.String(level.String()))
				attrs = append(attrs, logMessageKey.String(message))
				span.AddEvent("log", trace.WithAttributes(attrs...))
				if level >= zerolog.ErrorLevel {
					span.SetStatus(codes.Error, message)
				}
			}
		}
	}
}
