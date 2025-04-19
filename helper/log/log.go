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
var _logger ILogger

func init() {
	if _logger != nil {
		return
	}
	_defaultBuildFlow()
}

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

func Debug(msg string) {
	_logger.Debug(msg)
}
func Debugf(format string, args ...interface{}) {
	_logger.Debugf(format, args)
}
func Info(msg string) {
	_logger.Info(msg)
}
func Infof(format string, args ...interface{}) {
	_logger.Infof(format, args)
}
func Error(msg string) {
	_logger.Error(msg)
}
func Errorf(format string, args ...interface{}) {
	_logger.Errorf(format, args)
}
func Warn(msg string) {
	_logger.Error(msg)
}
func Warnf(format string, args ...interface{}) {
	_logger.Warnf(format, args)
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
	_logger = logger
	return
}

func (b *Builder) build() {
	if _logger != nil {
		return
	}
	var cores []zapcore.Core
	for _, sync := range b.syncs {
		cores = append(cores, zapcore.NewCore(b.enc, sync, b.level))
	}
	_logger = wrap(zap.New(zapcore.NewTee(cores...), b.options...))
	return
}

func GetLoggerInstance() ILogger {
	return logger
}

func QuickStart() {
	defaultBuildFlow()
}

func defaultBuilder() *Builder {
	return NewBuilder().
		LogTo(os.Stdout).
		LogWhen(zapcore.InfoLevel).
		EncodeWith(DefaultCfg, JsonFormat).
		EnableCaller().
		StackTraceOn(zapcore.ErrorLevel)
}

func defaultBuildFlow() {
	defaultBuilder().Build()
}

// _defaultBuildFlow 这个函数的意义时，当构建test单元时，可以快速使用日志器
// 而不需要重新写一段初始化日志的流程
func _defaultBuildFlow() {
	defaultBuilder().build()
}
