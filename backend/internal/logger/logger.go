package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	debugLog *log.Logger
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger
	fatalLog *log.Logger

	logFile    *os.File
	logLevel   LogLevel = INFO
	isInitialized bool
)

// Init 初始化日志系统
func Init(logDir string, level LogLevel) error {
	if isInitialized {
		return nil
	}

	// 创建日志目录
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// 创建日志文件
	logFilePath := filepath.Join(logDir, fmt.Sprintf("app_%s.log", time.Now().Format("2006-01-02")))
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}

	// 设置日志输出到文件和控制台
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// 初始化不同级别的日志记录器
	debugLog = log.New(multiWriter, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog = log.New(multiWriter, "[INFO]  ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLog = log.New(multiWriter, "[WARN]  ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog = log.New(multiWriter, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	fatalLog = log.New(multiWriter, "[FATAL] ", log.Ldate|log.Ltime|log.Lshortfile)

	logLevel = level
	isInitialized = true

	Info("Logger initialized successfully")
	return nil
}

// Close 关闭日志文件
func Close() error {
	if logFile != nil {
		return logFile.Close()
	}
	return nil
}

// Debug 调试日志
func Debug(format string, v ...interface{}) {
	if logLevel <= DEBUG && isInitialized {
		debugLog.Output(2, fmt.Sprintf(format, v...))
	}
}

// Info 信息日志
func Info(format string, v ...interface{}) {
	if logLevel <= INFO && isInitialized {
		infoLog.Output(2, fmt.Sprintf(format, v...))
	}
}

// Warn 警告日志
func Warn(format string, v ...interface{}) {
	if logLevel <= WARN && isInitialized {
		warnLog.Output(2, fmt.Sprintf(format, v...))
	}
}

// Error 错误日志
func Error(format string, v ...interface{}) {
	if logLevel <= ERROR && isInitialized {
		errorLog.Output(2, fmt.Sprintf(format, v...))
	}
}

// Fatal 致命错误日志(会导致程序退出)
func Fatal(format string, v ...interface{}) {
	if isInitialized {
		fatalLog.Output(2, fmt.Sprintf(format, v...))
	}
	os.Exit(1)
}

// SetLevel 设置日志级别
func SetLevel(level LogLevel) {
	logLevel = level
}

// GetLevel 获取当前日志级别
func GetLevel() LogLevel {
	return logLevel
}

// RotateLog 轮转日志文件(每天创建新文件)
func RotateLog(logDir string) error {
	if logFile != nil {
		logFile.Close()
	}

	logFilePath := filepath.Join(logDir, fmt.Sprintf("app_%s.log", time.Now().Format("2006-01-02")))
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to rotate log file: %v", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	debugLog.SetOutput(multiWriter)
	infoLog.SetOutput(multiWriter)
	warnLog.SetOutput(multiWriter)
	errorLog.SetOutput(multiWriter)
	fatalLog.SetOutput(multiWriter)

	Info("Log rotated successfully")
	return nil
}
