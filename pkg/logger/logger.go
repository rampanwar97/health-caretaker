package logger

import (
	"log"
	"os"
)

// Logger wraps the standard log package with additional functionality
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

// New creates a new logger instance
func New() *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info logs an info message
func (l *Logger) Info(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

// Error logs an error message
func (l *Logger) Error(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}

// Debug logs a debug message
func (l *Logger) Debug(format string, v ...interface{}) {
	if os.Getenv("DEBUG") == "true" {
		l.debugLogger.Printf(format, v...)
	}
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.errorLogger.Fatalf(format, v...)
}
