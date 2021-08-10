package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type (
	aboveDebug struct{}
	aboveInfo  struct{}
	aboveWarn  struct{}
)

var Logger *zap.Logger

func (l aboveDebug) Enabled(lv zapcore.Level) bool {
	return lv >= zapcore.DebugLevel
}

func (l aboveInfo) Enabled(lv zapcore.Level) bool {
	return lv >= zapcore.InfoLevel
}

func (l aboveWarn) Enabled(lv zapcore.Level) bool {
	return lv >= zapcore.WarnLevel
}

func makeDebugFilter() zapcore.LevelEnabler {
	return aboveDebug{}
}

func makeInfoFilter() zapcore.LevelEnabler {
	return aboveInfo{}
}

func makeErrorFilter() zapcore.LevelEnabler {
	return aboveWarn{}
}

// init
// 按照固定的模式初始化日志，调用方直接调用无需初始化
// 默认行为如下
// 存储路径:./log/
// 日志级别：各级别文件包含以上级别日志（DEBUG包含DEBUG/INFO/WARN/ERROR,INFO包含INFO/WARN/ERROR,ERROR包含WARN/ERROR）
// 其他：开启压缩，最大100MB，最多存50个文件，最久存28天
func init() {
	var encoder zapcore.Encoder
	encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	wDebug := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/debug.log",
		MaxSize:    500, // megabytes
		MaxBackups: 50,
		MaxAge:     28, // days
		Compress:   true,
	})
	wInfo := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/info.log",
		MaxSize:    300, // megabytes
		MaxBackups: 10,
		MaxAge:     28, // days
		Compress:   true,
	})
	wError := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/error.log",
		MaxSize:    100, // megabytes
		MaxBackups: 10,
		MaxAge:     28, // days
		Compress:   true,
	})

	coreDebug := zapcore.NewCore(
		encoder,
		wDebug,
		makeDebugFilter(),
	)

	coreInfo := zapcore.NewCore(
		encoder,
		wInfo,
		makeInfoFilter(),
	)
	coreError := zapcore.NewCore(
		encoder,
		wError,
		makeErrorFilter(),
	)

	Logger = zap.New(zapcore.NewTee(coreDebug, coreInfo, coreError), zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel))
}
