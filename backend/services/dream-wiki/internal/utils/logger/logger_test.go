package logger

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLoggerInitialization(t *testing.T) {
	logger1 := Init()
	logger2 := Init()

	if logger1 != logger2 {
		t.Error("Init() should return the same logger instance")
	}
}

func TestLogLevels(t *testing.T) {
	buf := bytes.Buffer{}
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	Info("test info", nil)
	Warn("test warn", nil)
	Error("test error", nil, nil)
	Infof("test infof: %d", 42)
	Warnf("test warnf: %s", "value")
	Errorf("test errorf: %v", "error")

	time.Sleep(100 * time.Millisecond)
	output := buf.String()

	expected := []string{
		"INFO",
		"test info",
		"WARN",
		"test warn",
		"ERROR",
		"test error",
		"test infof: 42",
		"test warnf: value",
		"test errorf: error",
	}

	for _, exp := range expected {
		if !strings.Contains(output, exp) {
			t.Errorf("Expected log output to contain %q", exp)
		}
	}
}

func TestLogWithFields(t *testing.T) {
	buf := bytes.Buffer{}
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	fields := map[string]any{
		"user_id": 123,
		"action":  "login",
	}
	Info("user action", fields)

	time.Sleep(100 * time.Millisecond)

	output := buf.String()

	expectedFields := []string{
		"user_id=123",
		"action=login",
	}

	for _, exp := range expectedFields {
		if !strings.Contains(output, exp) {
			t.Errorf("Expected log output to contain field %q", exp)
		}
	}
}

func TestErrorLogging(t *testing.T) {
	buf := bytes.Buffer{}
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	err := fmt.Errorf("database connection failed")
	Error("db operation", err, nil)

	time.Sleep(100 * time.Millisecond)

	output := buf.String()

	if !strings.Contains(output, "database connection failed") {
		t.Error("Expected error message in log output")
	}
}

func TestConcurrentLogging(t *testing.T) {
	buf := bytes.Buffer{}
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	var wg sync.WaitGroup
	count := 100

	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			Infof("concurrent log %d", n)
		}(i)
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	output := buf.String()

	for i := range count {
		if !strings.Contains(output, fmt.Sprintf("concurrent log %d", i)) {
			t.Errorf("Missing log message for iteration %d", i)
		}
	}
}

func TestLoggerClose(t *testing.T) {
	logger := Init()

	var closed bool
	go func() {
		logger.Close()
		closed = true
	}()

	time.Sleep(100 * time.Millisecond)

	if !closed {
		t.Error("Logger should be closed")
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Logging after close should not panic")
			}
		}()
		logger.log(LevelInfo, "test after close", nil)
	}()
}
