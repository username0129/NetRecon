package logger

import (
	"backend/internal/global"
	"backend/internal/util"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// customLevelEnabler 仅当等级严格匹配时启用日志记录
type customLevelEnabler zapcore.Level

// Enabled 实现了 zapcore.LevelEnabler 接口
func (cle customLevelEnabler) Enabled(level zapcore.Level) bool {
	return level == zapcore.Level(cle)
}

// InitializeLogger 初始化日志记录器
func InitializeLogger() *zap.Logger {
	if ok, _ := util.IsPathExist(global.Config.Zap.Director); !ok {
		if err := os.Mkdir(global.Config.Zap.Director, 0755); err != nil {
			fmt.Printf("无法创建日志目录: %v\n", err)
			os.Exit(1)
		}
	} // 创建日志目录

	cores := getZapCores()
	logger := zap.New(zapcore.NewTee(cores...))
	logger = logger.WithOptions(zap.AddCaller())
	return logger
}

func getZapCores() []zapcore.Core {
	cores := make([]zapcore.Core, 0, 7)
	for level := zapcore.DebugLevel; level <= zapcore.FatalLevel; level++ {
		cores = append(cores, getEncoderCore(level))
	}
	return cores
}

func getEncoderCore(l zapcore.Level) zapcore.Core {
	writer := getWriteSyncer(l.String())
	return zapcore.NewCore(getEncoder(), writer, customLevelEnabler(l))
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()              // 标准配置
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder          // 时间显示格式
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder         // 显示完整路径
	encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder // 小写带颜色

	switch global.Config.Zap.Format {
	case "json":
		return zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		return zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

func getWriteSyncer(level string) zapcore.WriteSyncer {
	fileWriter := NewRotate(level, "2006-01-02", global.Config.Zap.Director)
	if global.Config.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter))
	}
	return zapcore.AddSync(fileWriter)
}
