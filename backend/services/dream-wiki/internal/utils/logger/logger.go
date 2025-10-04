package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
	LevelFatal = "FATAL"
)

var (
	globalLogger *Logger
	initOnce     sync.Once
)

var osExit = os.Exit

type LogEntry struct {
	Level     string
	Message   string
	Timestamp time.Time
	Fields    map[string]any
}

type Logger struct {
	logChan  chan LogEntry
	stopChan chan struct{}
	wg       sync.WaitGroup
}

func Init() *Logger {
	initOnce.Do(func() {
		globalLogger = &Logger{
			logChan:  make(chan LogEntry, 100),
			stopChan: make(chan struct{}),
		}

		globalLogger.wg.Add(1)
		go globalLogger.processLogs()
	})
	return globalLogger
}

func Get() *Logger {
	if globalLogger == nil {
		return Init()
	}
	return globalLogger
}

func Close() {
	if globalLogger != nil {
		globalLogger.Close()
	}
}

func (l *Logger) processLogs() {
	defer l.wg.Done()

	for {
		select {
		case entry := <-l.logChan:
			logMsg := formatLogEntry(entry)
			log.Println(logMsg)

			if entry.Level == LevelFatal {
				osExit(1)
			}
		case <-l.stopChan:
			return
		}
	}
}

func Info(message string, fields map[string]any) {
	Get().log(LevelInfo, message, fields)
}

func Infof(format string, args ...any) {
	Get().log(LevelInfo, fmt.Sprintf(format, args...), nil)
}

func Warn(message string, fields map[string]any) {
	Get().log(LevelWarn, message, fields)
}

func Warnf(format string, args ...any) {
	Get().log(LevelWarn, fmt.Sprintf(format, args...), nil)
}

func Error(message string, err error, fields map[string]any) {
	if fields == nil {
		fields = make(map[string]any)
	}
	if err != nil {
		fields["error"] = err.Error()
	}
	Get().log(LevelError, message, fields)
}

func Errorf(format string, args ...any) {
	Get().log(LevelError, fmt.Sprintf(format, args...), nil)
}

func Fatal(message string, err error, fields map[string]any) {
	if fields == nil {
		fields = make(map[string]any)
	}
	if err != nil {
		fields["error"] = err.Error()
	}
	Get().log(LevelFatal, message, fields)
}

func Fatalf(format string, args ...any) {
	Get().log(LevelFatal, fmt.Sprintf(format, args...), nil)
}

func (l *Logger) log(level, message string, fields map[string]any) {
	if fields == nil {
		fields = make(map[string]any)
	}
	l.logChan <- LogEntry{
		Level:     level,
		Message:   message,
		Timestamp: time.Now(),
		Fields:    fields,
	}
}

func (l *Logger) Close() {
	close(l.stopChan)
	l.wg.Wait()
}

func formatLogEntry(entry LogEntry) string {
	fieldsStr := ""
	for k, v := range entry.Fields {
		fieldsStr += fmt.Sprintf("%s=%v ", k, v)
	}

	return fmt.Sprintf(
		"[%s] %s: %s %s",
		entry.Timestamp.Format(time.RFC3339),
		entry.Level,
		entry.Message,
		fieldsStr,
	)
}
