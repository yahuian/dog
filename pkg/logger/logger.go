package logger

import (
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type FileOption struct {
	// https://pkg.go.dev/gopkg.in/natefinch/lumberjack.v2?utm_source=godoc#Logger
	Filename string
	MaxSize  int  // 单个日志最大，单位： MB
	MaxAge   int  // 日志保留时间，单位：day
	Compress bool // .gz

	SplitTime int // 按时间切分，单位：hour，设置为 0 表示不按照时间切分
}

type myLogger struct {
	*zap.SugaredLogger
	zap.AtomicLevel
}

var instance *myLogger

// Init init logger if file is nil only output log info to os.Stdout
func Init(file *FileOption) {
	// console encoder and writer
	consoleEncoder, consoleWriter := getConsole()

	// json encoder and writer
	var (
		jsonEncoder zapcore.Encoder
		jsonWriter  zapcore.WriteSyncer
	)
	if file != nil {
		jsonEncoder, jsonWriter = getJSON(file)
	}

	// level
	atom := zap.NewAtomicLevel()

	// core
	cores := []zapcore.Core{
		zapcore.NewCore(consoleEncoder, consoleWriter, atom),
	}
	if file != nil {
		cores = append(cores, zapcore.NewCore(jsonEncoder, jsonWriter, atom))
	}
	core := zapcore.NewTee(cores...)

	// logger
	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	).Sugar()

	res := &myLogger{
		SugaredLogger: logger,
		AtomicLevel:   atom,
	}

	// set package global variable
	instance = res
}

func SetLevel(level string) {
	switch level {
	case "debug":
		instance.SetLevel(zap.DebugLevel)
	case "warn":
		instance.SetLevel(zap.WarnLevel)
	case "error":
		instance.SetLevel(zap.ErrorLevel)
	default:
		instance.SetLevel(zap.InfoLevel)
	}
}

func Sync() error {
	return instance.Sync()
}

func Debugln(args ...any) {
	instance.Debug(args...)
}

func Debugf(template string, args ...any) {
	instance.Debugf(template, args...)
}

func Infoln(args ...any) {
	instance.Info(args...)
}

func Infof(template string, args ...any) {
	instance.Infof(template, args...)
}

func Warnln(args ...any) {
	instance.Warn(args...)
}

func Warnf(template string, args ...any) {
	instance.Warnf(template, args...)
}

func Errorln(args ...any) {
	instance.Error(args...)
}

func Errorf(template string, args ...any) {
	instance.Errorf(template, args...)
}

func Panicln(args ...any) {
	instance.Panic(args...)
}

func Panicf(template string, args ...any) {
	instance.Panicf(template, args...)
}

func getConsole() (zapcore.Encoder, zapcore.WriteSyncer) {
	encoder := zapcore.NewConsoleEncoder(
		zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	)
	writer := zapcore.AddSync(os.Stdout)

	return encoder, writer
}

func getJSON(file *FileOption) (zapcore.Encoder, zapcore.WriteSyncer) {
	encoder := zapcore.NewJSONEncoder(
		zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	)

	// https://github.com/uber-go/zap/blob/master/FAQ.md#does-zap-support-log-rotation
	r := &lumberjack.Logger{
		Filename:  file.Filename,
		MaxSize:   file.MaxSize,
		MaxAge:    file.MaxAge,
		Compress:  file.Compress,
		LocalTime: true,
	}

	// https://github.com/natefinch/lumberjack/issues/17#issuecomment-185846531
	if file.SplitTime != 0 {
		go func() {
			for {
				<-time.After(time.Duration(file.SplitTime) * time.Hour)
				if err := r.Rotate(); err != nil {
					log.Printf("[dog] [ERROR] [logger] rotate err: %s", err.Error())
				}
			}
		}()
	}
	writer := zapcore.AddSync(r)

	return encoder, writer
}
