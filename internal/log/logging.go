package logger

import (
	"fmt"
	"github.com/gogoclouds/gogo/internal/conf"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"strings"
	"time"
)

var Log *zap.SugaredLogger

// 从配置文件映射结构

var (
	defLoggerPath           = "/logs"
	defLoggerFilenamePrefix = "log"
)

// Initialize logger handle
// appName log file prefix
func Initialize(appName string, conf conf.Log) {
	core := zapcore.NewCore( // 输出到日志文件
		setJSONEncoder(),
		setLoggerWriter(appName, conf),
		level2Int(conf.Level))

	outputConsole := zapcore.NewCore( // 输出到控制台
		setConsoleEncoder(),
		zapcore.Lock(os.Stdout),
		level2Int(conf.Level),
	)
	core = zapcore.NewTee(core, outputConsole)
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}

func setConsoleEncoder() zapcore.Encoder {
	ec := setEncoderConf()
	ec.EncodeLevel = zapcore.CapitalColorLevelEncoder // 终端输出 日志级别有颜色
	return zapcore.NewConsoleEncoder(ec)
}

func setJSONEncoder() zapcore.Encoder {
	ec := setEncoderConf()
	ec.EncodeLevel = zapcore.CapitalLevelEncoder // eg: info -> INFO
	return zapcore.NewConsoleEncoder(ec)
}

func setEncoderConf() zapcore.EncoderConfig {
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000") // 转换编码的时间戳
	return ec
}

// 可选 debug | info | error
// 默认 info
func level2Int(level string) zapcore.Level {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		log.Printf("No such level log: %s", level)
		return zapcore.InfoLevel
	}
}

func setLoggerWriter(appName string, conf conf.Log) zapcore.WriteSyncer {
	fName := makeFileName(conf.DirPath, appName)
	return zapcore.AddSync(
		&lumberjack.Logger{
			Filename:   fName,                 // 要写入的日志文件
			MaxSize:    int(conf.FileSizeMax), // 日志文件的大小（M）
			MaxBackups: 1,                     // 备份数量
			MaxAge:     int(conf.FileAgeMax),  // 存留天数
			Compress:   true,                  // 压缩
			LocalTime:  true,                  // 默认 UTC 时间
		})
}

// xxx/logs/xxx-service-2006-01-01-150405.log
func makeFileName(path, name string) string {
	if path = strings.TrimSpace(path); path == "" {
		path = defLoggerPath
	}
	if path = strings.TrimSpace(path); name == "" {
		name = defLoggerFilenamePrefix
	}
	nowTime := time.Now().Format("2006-01-02-150405")
	return fmt.Sprintf("%s/%s-%s.log", path, name, nowTime)
}
