package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
)

// atomicLevel 动态更新限制日志打印级别
var atomicLevel zap.AtomicLevel

func SetLevel(l LogLevel) {
	if l == "" {
		return
	}
	atomicLevel.UnmarshalText([]byte(l))
}

func InitZapLogger(conf Config) {
	atomicLevel = zap.NewAtomicLevel()
	go func() {
		for level := range conf.LevelCh {
			_ = atomicLevel.UnmarshalText([]byte(level))
		}
	}()
	fileCore := zapcore.NewCore( // 输出到日志文件
		setJSONEncoder(conf.TimeFormat, conf.FileJsonEncoder),
		setLoggerWriter(conf),
		atomicLevel,
	)
	consoleCore := zapcore.NewCore( // 输出到控制台
		setConsoleEncoder(conf.TimeFormat),
		zapcore.Lock(os.Stdout),
		atomicLevel,
	)

	core := zapcore.NewTee(fileCore, consoleCore)

	slogHandler := zapslog.NewHandler(core, zapslog.WithCaller(true), zapslog.AddStacktraceAt(16))
	InitSlog(slogHandler)

	// SetLogger(&ZapLogger{logger: zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))})
}

func setConsoleEncoder(timeFormat string) zapcore.Encoder {
	ec := setEncoderConf(timeFormat)
	ec.EncodeLevel = zapcore.CapitalColorLevelEncoder // 终端输出 日志级别有颜色
	return zapcore.NewConsoleEncoder(ec)
}

func setJSONEncoder(timeFormat string, isJsonEncoder bool) zapcore.Encoder {
	ec := setEncoderConf(timeFormat)
	ec.EncodeLevel = zapcore.CapitalLevelEncoder // eg: info -> INFO
	if isJsonEncoder {
		return zapcore.NewJSONEncoder(ec)
	}
	return zapcore.NewConsoleEncoder(ec)
}

func setEncoderConf(timeFmt string) zapcore.EncoderConfig {
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeTime = zapcore.TimeEncoderOfLayout(timeFmt) // 转换编码的时间戳
	return ec
}

func setLoggerWriter(conf Config) zapcore.WriteSyncer {
	fName := conf.makeFilename()
	return zapcore.AddSync(
		&lumberjack.Logger{
			Filename:   fName,                 // 要写入的日志文件
			MaxSize:    int(conf.FileMaxSize), // 日志文件的大小（M）
			MaxAge:     int(conf.FileMaxAge),  // 存留天数
			MaxBackups: 1,                     // 备份数量
			Compress:   conf.FileCompress,     // 压缩
			LocalTime:  true,                  // 默认 UTC 时间
		})
}