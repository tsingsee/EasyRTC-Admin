// Copyright 2020 TSINGSEE.
// http://www.tsingsee.com
// 日志模块配置
// Creat By Sam
// History (Name, Time, Desc)
// (Sam, 20200506, 创建文件)
package util

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
    "os"
    "time"
)

// 操作日志
var operationLogger *zap.Logger

// 获取日志
func GetLogger() *zap.Logger {
    if operationLogger == nil {
        operationLogger = createLog("./easyrtc.log", true)
    }

    return operationLogger
}

// 初始化 Log 模块
func createLog(filePath string, develop bool) *zap.Logger {
    //logpath := filepath.Join(GetConf().DirLogs, fileName)

    hook := lumberjack.Logger{
        Filename:   filePath, // 日志文件路径
        MaxSize:    10,       // 每个日志文件保存的最大尺寸 单位：M
        MaxBackups: 10,       // 日志文件最多保存多少个备份
        MaxAge:     30,       // 文件最多保存多少天
        Compress:   false,    // 是否压缩
    }

    encoderConfig := zapcore.EncoderConfig{
        TimeKey:       "time",
        LevelKey:      "level",
        NameKey:       "monitorLogger",
        CallerKey:     "linenum",
        MessageKey:    "msg",
        StacktraceKey: "stacktrace",
        LineEnding:    zapcore.DefaultLineEnding,
        EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
        EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
            enc.AppendString(t.Format("2006-01-02 15:04:05"))
        },
        EncodeDuration: zapcore.SecondsDurationEncoder, //
        EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
        EncodeName:     zapcore.FullNameEncoder,
    }

    // 设置日志级别
    atomicLevel := zap.NewAtomicLevel()
    atomicLevel.SetLevel(zap.InfoLevel)

    core := zapcore.NewCore(
        zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
        zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
        atomicLevel,                                                                     // 日志级别
    )

    /*// 设置初始化字段
    filed := zap.Fields(zap.String("monitor", "dssfolder"))
    */

    var logger *zap.Logger
    if develop {
        // 开启文件及行号
        // 开启开发模式，堆栈跟踪
        caller := zap.AddCaller()
        develop := zap.Development()
        logger = zap.New(core, caller, develop)
    } else {
        logger = zap.New(core)
    }

    return logger
}
