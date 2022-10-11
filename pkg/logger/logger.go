package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

const (
	DebugFileName = "log/info.log"
	ErrorFileName = "log/error.log"
	TimeFormat    = "2006-01"
)

var Log *zap.SugaredLogger

func init() {
	encoder := getEncoder()

	dFile := fmt.Sprintf("%s.%s", DebugFileName, time.Now().Format(TimeFormat))
	eFile := fmt.Sprintf("%s.%s", ErrorFileName, time.Now().Format(TimeFormat))

	logInfo := zapcore.NewCore(encoder, getWriter(dFile), zapcore.InfoLevel)
	logError := zapcore.NewCore(encoder, getWriter(eFile), zapcore.ErrorLevel)

	logDebug := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)

	core := zapcore.NewTee(logError, logInfo, logDebug)

	Log = zap.New(core, zap.AddCaller()).Sugar()
	defer Log.Sync()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "time",
		StacktraceKey: "stacktrace",
		EncodeTime: func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString("[" + t.Format("2006-01-02 15:04:05.000") + "]")
		},

		CallerKey: "file",
		EncodeCaller: func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString("[FilePath] " + caller.TrimmedPath() + " ")
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,

		EncodeLevel: ColorEncoder,
		LineEnding:  zapcore.DefaultLineEnding,
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getWriter(fileName string) zapcore.WriteSyncer {
	ll := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    10,
		MaxAge:     30,
		MaxBackups: 5,
		Compress:   false,
	}
	ws := io.MultiWriter(ll, os.Stdout)
	return zapcore.AddSync(ws)
}
