package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	zl "github.com/rs/zerolog/log"
)

type Logger interface {
	Debug(message string, args ...any)
	Info(message string, args ...any)
	Warn(message string, args ...any)
	Error(message string, args ...any)
	Fatal(message string, args ...any)
}

func New(logLevel, serviceName string) (Logger, error) {
	return newZerolog(logLevel, serviceName)
}

type logger struct {
	inner *zerolog.Logger
}

func newZerolog(logLevel string, serviceName string) (Logger, error) {
	var l zerolog.Level
	out := io.Writer(os.Stderr)

	switch logLevel {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
		out = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.DateTime,
			NoColor:    false,
		}
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	zl.Logger = zerolog.New(out).
		With().
		Str("service", serviceName).
		Timestamp().
		CallerWithSkipFrameCount(4).
		Logger()

	return &logger{
		inner: &zl.Logger,
	}, nil
}

// Debug -.
func (l *logger) Debug(message string, args ...any) {
	l.log(zerolog.DebugLevel, message, args...)
}

// Info -.
func (l *logger) Info(message string, args ...any) {
	l.log(zerolog.InfoLevel, message, args...)
}

// Warn -.
func (l *logger) Warn(message string, args ...any) {
	l.log(zerolog.WarnLevel, message, args...)
}

// Error -.
func (l *logger) Error(message string, args ...any) {
	l.log(zerolog.ErrorLevel, message, args...)
}

// Fatal -.
func (l *logger) Fatal(message string, args ...any) {
	l.log(zerolog.FatalLevel, message, args...)

	os.Exit(1)
}

func (l *logger) log(level zerolog.Level, message string, args ...any) {
	if len(args) == 0 {
		l.inner.WithLevel(level).Msg(message)
		return
	}

	event := l.inner.WithLevel(level)

	// Process paired arguments
	pairsCount := len(args) / 2
	for i := 0; i < pairsCount; i++ {
		idx := i * 2
		var key string
		if k, ok := args[idx].(string); ok {
			key = k
		} else {
			key = fmt.Sprintf("arg%d", idx)
		}

		value := args[idx+1]

		switch v := value.(type) {
		case string:
			event = event.Str(key, v)
		case int:
			event = event.Int(key, v)
		case error:
			event = event.Err(v)
		default:
			event = event.Interface(key, v)
		}
	}

	// If there's an unpaired last element
	if len(args)%2 != 0 {
		lastIdx := len(args) - 1
		key := fmt.Sprintf("arg%d", lastIdx)
		event = event.Interface(key, args[lastIdx])
	}

	event.Msg(message)
}
