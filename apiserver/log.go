package apiserver

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
)

func NewRotateWriter() io.Writer {
	return &lumberjack.Logger{
		Filename:   "logger.log",
		MaxSize:    10, // 单位：MB
		MaxBackups: 5,
		MaxAge:     7, // 单位：天
		Compress:   true,
	}
}
