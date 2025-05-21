package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

var log *zap.Logger

func NewSlogCore(s *SlogConfig) (*zap.Logger, error) {
	if s.Mode == "dev" {
		logger, _ := zap.NewDevelopment() // zap.ReplaceGlobals(logger)
		log = logger
	} else {
		fileDir := filepath.Dir(s.Filename)
		if _, err := os.Stat(fileDir); os.IsNotExist(err) {
			if err = os.MkdirAll(fileDir, 0770); err != nil {
				fmt.Printf("Failed to create data directory: %s, err: %+v\n", fileDir, err)
				return nil, err
			}
		}
		lumberJackLogger := &lumberjack.Logger{
			Filename:   s.Filename,
			MaxSize:    s.MaxSize,
			MaxAge:     s.MaxAge,
			MaxBackups: s.MaxBackups,
			Compress:   s.Compress,
		}

		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		consoleCore := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), selectLevel(s.ProdLevel))

		encoder := GetEncoder()
		syncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger))
		fileCore := zapcore.NewCore(encoder, syncer, zapcore.InfoLevel)

		core := zapcore.NewTee(fileCore, consoleCore)

		log = zap.New(core, zap.AddCaller())
		defer log.Sync()

	}
	return log, nil
}

func selectLevel(level string) zapcore.Level {
	switch level {
	case "info":
		return zapcore.InfoLevel
	case "debug":
		return zapcore.DebugLevel
	case "error":
		return zapcore.ErrorLevel
	case "warning":
		return zapcore.WarnLevel
	default:
		return zapcore.InfoLevel
	}
}

func GetEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
