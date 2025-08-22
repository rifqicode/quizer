package utils

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

// LogError logs an error with additional context
func LogError(operation string, err error, context ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	log.Printf("[ERROR] %s in %s (%s:%d) - %v - Context: %v",
		operation, funcName, file, line, err, context)
}

// LogInfo logs informational messages
func LogInfo(message string, context ...interface{}) {
	log.Printf("[INFO] %s - Context: %v", message, context)
}

// LogWarning logs warning messages
func LogWarning(message string, context ...interface{}) {
	log.Printf("[WARNING] %s - Context: %v", message, context)
}

// SafeExecute executes a function with panic recovery
func SafeExecute(operation string, fn func() error) error {
	defer func() {
		if r := recover(); r != nil {
			LogError(operation, fmt.Errorf("panic recovered: %v", r))
		}
	}()

	return fn()
}

// Retry executes a function with retry logic
func Retry(operation string, attempts int, delay time.Duration, fn func() error) error {
	var lastErr error

	for i := 0; i < attempts; i++ {
		if err := fn(); err != nil {
			lastErr = err
			LogWarning(fmt.Sprintf("%s attempt %d failed", operation, i+1), err)

			if i < attempts-1 {
				time.Sleep(delay)
			}
		} else {
			return nil
		}
	}

	return fmt.Errorf("%s failed after %d attempts: %w", operation, attempts, lastErr)
}
