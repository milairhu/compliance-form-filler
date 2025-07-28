package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/urfave/cli/v3"
)

type Logger struct {
	logger *zerolog.Logger
}

// DefaultLogger is the default logger for the package level functions.
// It is initialized with the default configuration.
var DefaultLogger *Logger

// Log format enum
type LogFormat string

const (
	// TextFormat is the default log format
	TextFormat LogFormat = "text"
	// JSONFormat is the log format in JSON
	JSONFormat LogFormat = "json"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logger := zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		Logger()

	DefaultLogger = &Logger{logger: &logger}
}

func NewFromCliContext(c *cli.Command) *Logger {
	verbose := c.Bool("verbose")
	logFormat := c.String("log-format")

	DefaultLogger = New(verbose, LogFormat(logFormat))
	DefaultLogger = DefaultLogger.WithBaseFields(uuid.New().String())

	if orgID := c.String("organization-id"); orgID != "" {
		DefaultLogger = DefaultLogger.WithField("organization-id", orgID)
		DefaultLogger = DefaultLogger.WithField("workspace-id", c.String("workspace-id"))
		DefaultLogger = DefaultLogger.WithField("instance-id", c.String("instance-id"))
	}

	return DefaultLogger
}

// New creates a new logger with the specified configuration.
func New(isDebug bool, format LogFormat) *Logger {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.DebugLevel
	}

	var logWriter io.Writer = os.Stdout
	if format == TextFormat {
		logWriter = zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	}

	zerolog.SetGlobalLevel(logLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logger := zerolog.
		New(logWriter).
		With().
		Timestamp().
		Logger()

	return &Logger{logger: &logger}
}

// Nop function returns a logger that does nothing (Nop).
func Nop() *Logger {
	logger := zerolog.Nop()

	return &Logger{logger: &logger}
}

// NewConsole creates a new logger that writes to the console.
func NewConsole(isDebug bool) *Logger {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return &Logger{logger: &logger}
}

// Output duplicates the global logger and sets w as its output.
func (l *Logger) Output(w io.Writer) zerolog.Logger {
	return l.logger.Output(w)
}

// With creates a child logger with the field added to its context.
func (l *Logger) With() zerolog.Context {
	return l.logger.With()
}

// Level creates a child logger with the minimum accepted level set to level.
func (l *Logger) Level(level zerolog.Level) zerolog.Logger {
	return l.logger.Level(level)
}

// Sample returns a logger with the s sampler.
func (l *Logger) Sample(s zerolog.Sampler) zerolog.Logger {
	return l.logger.Sample(s)
}

// Hook returns a logger with the h Hook.
func (l *Logger) Hook(h zerolog.Hook) zerolog.Logger {
	return l.logger.Hook(h)
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Debug() *zerolog.Event {
	return l.logger.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Info() *zerolog.Event {
	return l.logger.Info()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Warn() *zerolog.Event {
	return l.logger.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Error() *zerolog.Event {
	return l.logger.Error()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Panic() *zerolog.Event {
	return l.logger.Panic()
}

// WithLevel starts a new message with level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) WithLevel(level zerolog.Level) *zerolog.Event {
	return l.logger.WithLevel(level)
}

// Log starts a new message with no level. Setting zerolog.GlobalLevel to
// zerolog.Disabled will still disable events produced by this method.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Log() *zerolog.Event {
	return l.logger.Log()
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(v ...interface{}) {
	l.logger.Print(v...)
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

func (l *Logger) WithContext(ctx context.Context) context.Context {
	return l.logger.WithContext(ctx)
}

// Ctx returns the Logger associated with the ctx. If no logger
// is associated, a disabled logger is returned.
func Ctx(ctx context.Context) *Logger {
	return &Logger{logger: zerolog.Ctx(ctx)}
}

// UpdateContext updates the logger context with the given function.
func (l *Logger) UpdateContext(update func(c zerolog.Context) zerolog.Context) {
	l.logger.UpdateContext(update)
}

// WithBaseFields function is adding fields, like toolbox-run-id
func (l *Logger) WithBaseFields(toolboxRunID string) *Logger {
	logger := DefaultLogger.
		With().
		Str("toolbox-run-id", toolboxRunID).
		Logger()

	return &Logger{logger: &logger}
}

func (l *Logger) WithField(key string, value string) *Logger {
	logger := DefaultLogger.
		With().
		Str(key, value).
		Logger()

	return &Logger{logger: &logger}
}

func (l *Logger) GetZeroLogger() *zerolog.Logger {
	return l.logger
}
