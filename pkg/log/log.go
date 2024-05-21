package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

const (
	develop = "develop"
)

type Config struct {
	Path       string `json:"path"`
	FileName   string `json:"file_name"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	Mode       string `json:"mode"`
	Encoder    string `json:"encoder"`
}

func configure(config Config) zapcore.WriteSyncer {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.Path + "/" + config.FileName + ".log",
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     30,
		Compress:   true,
	})
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stderr),
		zapcore.AddSync(w),
	)
}

func GetLogger(config Config) *zap.SugaredLogger {
	var logLevel = zapcore.InfoLevel
	if config.Mode == develop {
		logLevel = zapcore.DebugLevel
	}

	var encoderCfg zapcore.EncoderConfig
	if config.Mode == develop {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey = "TIME"
	encoderCfg.NameKey = "NAME"
	encoderCfg.MessageKey = "MESSAGE"
	encoderCfg.EncodeDuration = zapcore.NanosDurationEncoder
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.FunctionKey = "FUNC"
	var encoder zapcore.Encoder
	if config.Encoder == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}
	core := zapcore.NewCore(encoder, configure(config), zap.NewAtomicLevelAt(logLevel))
	loggerZap := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(0))
	return loggerZap.Sugar()
}
