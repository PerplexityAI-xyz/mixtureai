package log

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mixtureai/config"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var errorLogger *zap.SugaredLogger

func getWriter(filename string) io.Writer {
	hook, err := rotatelogs.New(
		strings.Replace(filename, ".log", "", -1) + "-%Y%m%d%H.log", // 没有使用go风格反人类的format格式
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func init() {
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	infoWriter := getWriter(fmt.Sprintf("%s/info.log", config.C.LogDir))
	errorWriter := getWriter(fmt.Sprintf("%s/error.log", config.C.LogDir))

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	)

	log := zap.New(core, zap.AddCaller()) // 需要传入 zap.AddCaller() 才会显示打日志点的文件名和行数, 有点小坑
	errorLogger = log.Sugar()
}

func Sync() {
	errorLogger.Sync()
}

func Debug(args ...interface{}) {
	errorLogger.Debug(args...)
}
func Debugf(template string, args ...interface{}) {
	errorLogger.Debugf(template, args...)
}
func Info(args ...interface{}) {
	errorLogger.Info(args...)
}
func Infof(template string, args ...interface{}) {
	errorLogger.Infof(template, args...)
}
func Warn(args ...interface{}) {
	errorLogger.Warn(args...)
}
func Warnf(template string, args ...interface{}) {
	errorLogger.Warnf(template, args...)
}
func Error(args ...interface{}) {
	errorLogger.Error(args...)
}
func Errorf(template string, args ...interface{}) {
	errorLogger.Errorf(template, args...)
}
func DPanic(args ...interface{}) {
	errorLogger.DPanic(args...)
}
func DPanicf(template string, args ...interface{}) {
	errorLogger.DPanicf(template, args...)
}
func Panic(args ...interface{}) {
	errorLogger.Panic(args...)
}
func Panicf(template string, args ...interface{}) {
	errorLogger.Panicf(template, args...)
}
func Fatal(args ...interface{}) {
	errorLogger.Fatal(args...)
}
func Fatalf(template string, args ...interface{}) {
	errorLogger.Fatalf(template, args...)
}

func Indentf(template string, args ...interface{}) {
	b, err := json.MarshalIndent(args, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	errorLogger.Infof(template, string(b))
}

type RestyLogger struct {
}

func (l *RestyLogger) Errorf(format string, v ...interface{}) {
	errorLogger.Errorf("ERROR RESTY "+format, v...)
}

func (l *RestyLogger) Warnf(format string, v ...interface{}) {
	errorLogger.Warnf("WARN RESTY "+format, v...)
}

func (l *RestyLogger) Debugf(format string, v ...interface{}) {
	errorLogger.Debugf("DEBUG RESTY "+format, v...)
}

type ESLogger struct {
}

func (l ESLogger) Printf(format string, v ...interface{}) {
	errorLogger.Infof("[ES] "+format, v...)
}

type CronWriter struct {
}

func (w CronWriter) Write(p []byte) (n int, err error) {
	n = len(p)

	return
}

type DBWriter struct {
}

func (w DBWriter) Printf(format string, args ...interface{}) {

	errorLogger.Infof("[DB]"+format, args...)
}

type FiberWriter struct {
}

func (w FiberWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	errorLogger.Infof("[FIBER] %v ", string(p))
	return
}
