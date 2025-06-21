package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
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
	builtinLog()
}

func Init() {
	NewBuilder().
		LogTo(os.Stdout).
		LogWhen(zapcore.InfoLevel).
		EncodeWith(DefaultCfg, JsonFormat).
		EnableCaller().
		StackTraceOn(zapcore.ErrorLevel)
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
	Fatal(msg string)
	Fatalf(format string, args ...interface{})
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
	l.logger.Debug(fmt.Sprintf(format, args...))
}

func (l *Log) Info(msg string) {
	l.logger.Info(msg)
}

func (l *Log) Infof(format string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *Log) Error(msg string) {
	l.logger.Error(msg)
}

func (l *Log) Errorf(format string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

func (l *Log) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l *Log) Warnf(format string, args ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

func (l *Log) Fatal(msg string) {
	l.logger.Fatal(msg)
}

func (l *Log) Fatalf(format string, args ...interface{}) {
	l.logger.Fatal(fmt.Sprintf(format, args...))
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

func (b *Builder) LogTo(writer ...io.Writer) *Builder {
	for _, w := range writer {
		b.syncs = append(b.syncs, zapcore.AddSync(w))
	}
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

func defaultBuilder() *Builder {
	return NewBuilder().
		LogTo(os.Stdout).
		LogWhen(zapcore.InfoLevel).
		EncodeWith(DefaultCfg, JsonFormat).
		EnableCaller().
		StackTraceOn(zapcore.ErrorLevel)
}

// builtinLog 这个函数的意义时，当构建test单元时，可以快速使用日志器
// 而不需要重新写一段初始化日志的流程
func builtinLog() {
	defaultBuilder().build()
}

func AsZapLoggerPlugin() *zap.Logger {
	return _logger.(*Log).logger
}

type GormLoggerPlugin struct {
	Logger ILogger
}

const (
	infoStr      = "%s\n[info] "
	warnStr      = "%s\n[warn] "
	errStr       = "%s\n[error] "
	traceStr     = "%s\n[%.3fms] [rows:%v] %s"
	traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
	traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
)

func (g *GormLoggerPlugin) Printf(format string, v ...interface{}) {
	if strings.HasPrefix(format, infoStr) {
		g.Logger.Infof(format, v[1:]...)
	} else if strings.HasPrefix(format, warnStr) {
		g.Logger.Warnf(format, v[1:]...)
	} else if strings.HasPrefix(format, errStr) {
		g.Logger.Errorf(format, v[1:]...)
	} else if strings.HasPrefix(format, traceStr) {
		g.Logger.Debugf(format, v[1:]...)
	} else if strings.HasPrefix(format, traceWarnStr) {
		g.Logger.Warnf(format, v[1:]...)
	} else if strings.HasPrefix(format, traceErrStr) {
		g.Logger.Errorf(format, v[1:]...)
	} else {
		g.Logger.Debugf(format, v...)
	}
}

func AsGormLoggerPlugin() *GormLoggerPlugin {
	return &GormLoggerPlugin{
		Logger: _logger,
	}
}
