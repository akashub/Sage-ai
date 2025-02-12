// backend/go/pkg/logger/logger.go
package logger

import (
    "fmt"
    "log"
    "os"
    "time"
)

var (
    InfoLogger  *log.Logger
    ErrorLogger *log.Logger
    DebugLogger *log.Logger
)

func init() {
    // Create logs directory if it doesn't exist
    err := os.MkdirAll("logs", 0755)
    if err != nil {
        log.Fatal("Failed to create logs directory")
    }

    // Create or append to log file
    currentTime := time.Now().Format("2006-01-02")
    logFile, err := os.OpenFile(fmt.Sprintf("logs/sage_ai_%s.log", currentTime), 
        os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal("Failed to open log file")
    }

    // Initialize loggers
    InfoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    ErrorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
    DebugLogger = log.New(logFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

    // Also write to stdout
    InfoLogger.SetOutput(os.Stdout)
    ErrorLogger.SetOutput(os.Stdout)
    DebugLogger.SetOutput(os.Stdout)
}