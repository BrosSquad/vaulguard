package log

import (
	"context"
	"github.com/rs/zerolog"
	"gorm.io/gorm/logger"
	"io"
	"os"
	"time"
)

type (
	DBLevel int
	Logger  struct {
		addConsole bool
		level      zerolog.Level
		logger     zerolog.Logger
		writers    []io.Writer
	}
)

var Log *Logger

const (
	DbSilent DBLevel = iota + 1
	DbError
	DbWarn
	DbInfo
)

func GetLogLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "error":
		return zerolog.ErrorLevel
	case "warn":
		return zerolog.WarnLevel
	case "fatal":
		return zerolog.FatalLevel
	case "silent":
		return zerolog.NoLevel
	default:
		return zerolog.InfoLevel
	}
}

func GetDbLogLevel(level string) DBLevel {
	switch level {
	case "info":
		fallthrough
	case "debug":
		return DbInfo
	case "error":
		return DbError
	case "warn":
		return DbWarn
	case "silent":
		return DbSilent
	default:
		return DbError
	}
}

func NewVaulGuardLogger(level zerolog.Level, addConsole bool, writers ...io.Writer) *Logger {
	if addConsole {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout})
	}
	zerolog.SetGlobalLevel(level)

	return &Logger{
		addConsole: addConsole,
		level:      level,
		writers:    writers,
		logger:     zerolog.New(zerolog.MultiLevelWriter(writers...)).With().Timestamp().Logger(),
	}
}

func SetDefaultLogger(logger *Logger) {
	Log = logger
}

func (l Logger) LogMode(level logger.LogLevel) logger.Interface {
	zeroLogLevel := zerolog.NoLevel

	switch level {
	case logger.Warn:
		zeroLogLevel = zerolog.WarnLevel
	case logger.Error:
		zeroLogLevel = zerolog.ErrorLevel
	case logger.Info:
		zeroLogLevel = zerolog.InfoLevel
	case logger.Silent:
		zeroLogLevel = zerolog.NoLevel
	}

	return NewVaulGuardLogger(zeroLogLevel, l.addConsole, l.writers...)
}

func (l Logger) Debug(format string, data ...interface{}) {
	l.logger.Debug().Msgf(format, data...)
}

func (l Logger) Info(ctx context.Context, format string, data ...interface{}) {
	l.logger.WithContext(ctx)
	l.logger.Info().Msgf(format, data...)
}

func (l Logger) Warn(ctx context.Context, format string, data ...interface{}) {
	l.logger.WithContext(ctx)
	l.logger.Warn().Msgf(format, data...)
}

func (l Logger) Error(ctx context.Context, format string, data ...interface{}) {
	l.logger.WithContext(ctx)
	l.logger.Error().Msgf(format, data...)
}

func (l Logger) Errorf(err error, format string, data ...interface{}) {
	l.logger.Error().Err(err).Msgf(format, data...)
}

func (l Logger) Fatalf(err error, format string, data ...interface{}) {
	l.logger.Fatal().Err(err).Msgf(format, data...)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, callback func() (string, int64), err error) {
	message, number := callback()
	l.logger.WithContext(ctx)
	l.logger.
		Error().
		Err(err).
		Int64("Gorm: ", number).
		Time("Begin: ", begin).
		Msg(message)
}

func (l Logger) Printf(format string, data ...interface{}) {
	l.logger.WithLevel(l.level).Msgf(format, data...)
}
