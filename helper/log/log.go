package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

const (
	JsonFormat    = "JSON"
	ConsoleFormat = "CONSOLE"
)

var DefaultCfg = zapcore.EncoderConfig{
	MessageKey:       "msg",
	LevelKey:         "level",
	TimeKey:          "time",
	NameKey:          "name",
	CallerKey:        "caller",
	StacktraceKey:    "stackTrace",
	LineEnding:       "\r\n",
	EncodeTime:       zapcore.RFC3339TimeEncoder,
	EncodeLevel:      zapcore.CapitalLevelEncoder,
	EncodeCaller:     zapcore.ShortCallerEncoder,
	ConsoleSeparator: "\t",
}

var logger ILogger

type ILogger interface {
	Debug(msg string)
	Debugf(format string, args ...interface{})
	Info(msg string)
	Infof(format string, args ...interface{})
	Error(msg string)
	Errorf(format string, args ...interface{})
	Warn(msg string)
	Warnf(format string, args ...interface{})
}

type Log struct {
	logger *zap.Logger
}

func (l *Log) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l *Log) Debugf(format string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(format, args))
}

func (l *Log) Info(msg string) {
	l.logger.Info(msg)
}

func (l *Log) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args))
}

func (l *Log) Error(msg string) {
	l.logger.Error(msg)
}

func (l *Log) Errorf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args))
}

func (l *Log) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l *Log) Warnf(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args))
}

func wrap(logger *zap.Logger) ILogger {
	return &Log{
		logger: logger,
	}
}

type Builder struct {
	enc   zapcore.Encoder
	syncs []zapcore.WriteSyncer
	level zapcore.Level

	options []zap.Option
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) EncodeWith(cfg zapcore.EncoderConfig, format string) *Builder {
	switch format {
	case ConsoleFormat:
		b.enc = zapcore.NewConsoleEncoder(cfg)
	default:
		b.enc = zapcore.NewJSONEncoder(cfg)
	}
	return b
}

func (b *Builder) LogTo(w io.Writer) *Builder {
	b.syncs = append(b.syncs, zapcore.AddSync(w))
	return b
}

func (b *Builder) LogWhen(level zapcore.Level) *Builder {
	b.level = level
	return b
}

func (b *Builder) EnableCaller() *Builder {
	b.options = append(b.options, zap.AddCaller())
	b.options = append(b.options, zap.AddCallerSkip(0))
	return b
}

func (b *Builder) StackTraceOn(level zapcore.Level) *Builder {
	b.options = append(b.options, zap.AddStacktrace(level))
	return b
}

func (b *Builder) Build() {
	if logger != nil {
		return
	}
	var cores []zapcore.Core
	for _, sync := range b.syncs {
		cores = append(cores, zapcore.NewCore(b.enc, sync, b.level))
	}
	logger = wrap(zap.New(zapcore.NewTee(cores...), b.options...))
	return
}

func GetLoggerInstance() ILogger {
	return logger
}

func QuickStart() {
	NewBuilder().
		LogTo(os.Stdout).
		LogWhen(zapcore.InfoLevel).
		EncodeWith(DefaultCfg, JsonFormat).
		EnableCaller().
		StackTraceOn(zapcore.ErrorLevel).
		Build()
}
