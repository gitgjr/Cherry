package zlog

import (
	"io"
	"os"
	"path"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 简单封装一下对 zap 日志库的使用
// 使用方式：
// zlog.Debug("hello", zap.String("name", "Kevin"), zap.Any("arbitraryObj", dummyObject))
// zlog.Info("hello", zap.String("name", "Kevin"), zap.Any("arbitraryObj", dummyObject))
// zlog.Warn("hello", zap.String("name", "Kevin"), zap.Any("arbitraryObj", dummyObject))

var logger *zap.Logger

func init() {
	// 日志Encoder 还是JSONEncoder，把日志行格式化成JSON格式的
	encoder := getEncoder()

	fileWriteSyncer := getLogWriter()

	core := zapcore.NewCore(encoder, fileWriteSyncer, zapcore.DebugLevel)

	logger = zap.New(core)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.OpenFile("/test.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	// 利用io.MultiWriter支持文件和终端两个输出目标
	ws := io.MultiWriter(file, os.Stdout)
	return zapcore.AddSync(ws)
}

// func getFileLogWriter() (writeSyncer zapcore.WriteSyncer) {
// 	// 使用 lumberjack 实现 log rotate
// 	lumberJackLogger := &lumberjack.Logger{
// 		Filename:   "/tmp/test.log",
// 		MaxSize:    100, // 单个文件最大100M
// 		MaxBackups: 60,  // 多于 60 个日志文件后，清理较旧的日志
// 		MaxAge:     1,   // 一天一切割
// 		Compress:   false,
// 	}

// 	return zapcore.AddSync(lumberJackLogger)
// }

func Info(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Error(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	callerFields := getCallerInfoForLog()
	fields = append(fields, callerFields...)
	logger.Warn(message, fields...)
}

func getCallerInfoForLog() (callerFields []zap.Field) {

	pc, file, line, ok := runtime.Caller(2) // 回溯两层，拿到写日志的调用方的函数信息
	if !ok {
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	funcName = path.Base(funcName) //Base函数返回路径的最后一个元素，只保留函数名

	callerFields = append(callerFields, zap.String("func", funcName), zap.String("file", file), zap.Int("line", line))
	return
}
