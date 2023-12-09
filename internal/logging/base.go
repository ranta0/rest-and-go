package logging

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm/logger"
)

type LoggerFile struct {
	logLevel logger.LogLevel
	logger   *log.Logger
	file     *os.File
}

func NewLoggerFile(dir string, filename string) (*LoggerFile, error) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	logFilePath := filepath.Join(dir, filename)
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	mw := io.MultiWriter(os.Stdout, file)
	logger := log.New(mw, "", log.LstdFlags)

	loggerFile := &LoggerFile{
		logger: logger,
		file: file,
	}

	loggerFile.SetProd()
	return loggerFile, nil
}

// Implementation of the default Logger interface for gorm
func (l *LoggerFile) LogMode(level logger.LogLevel) logger.Interface {
	l.logLevel = level
	return l
}

func (l *LoggerFile) Info(ctx context.Context, msg string, data ...interface{}) {
	l.writeLogFile("info", logger.Info, msg, data...)
}

func (l *LoggerFile) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.writeLogFile("warn", logger.Warn, msg, data...)
}

func (l *LoggerFile) Error(ctx context.Context, msg string, data ...interface{}) {
	l.writeLogFile("error", logger.Error, msg, data...)
}

func (l *LoggerFile) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	l.writeLogFile("trace", logger.Silent, fmt.Sprintf("[%.3fms] %s", float64(time.Since(begin).Microseconds())/1000, sql), "Rows affected:", rows, "Error:", err)
}

func (l *LoggerFile) writeLogFile(levelType string, logLevel logger.LogLevel, msg string, data ...interface{}) {
	if logLevel >= l.logLevel {
		logMessage := fmt.Sprintf("[%s] %s %v\n", levelType, msg, data)
		l.file.WriteString(logMessage)
	}
}

func (l *LoggerFile) GetLogger() *log.Logger {
	return l.logger
}

func (l *LoggerFile) SetProd() {
	l.logLevel = logger.Error
}

func (l *LoggerFile) SetDebug() {
	l.logLevel = logger.Silent
}

func (l *LoggerFile) Println(v ...interface{}) {
	l.logger.Println(v...)
}

func (l *LoggerFile) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

func (l *LoggerFile) Close() {
	l.file.Close()
}
