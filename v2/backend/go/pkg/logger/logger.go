// // backend/go/pkg/logger/logger.go
// package logger

// import (
//     "fmt"
//     "log"
//     "os"
//     "time"
// )

// var (
//     InfoLogger  *log.Logger
//     ErrorLogger *log.Logger
//     DebugLogger *log.Logger
// )

// func init() {
//     // Create logs directory if it doesn't exist
//     err := os.MkdirAll("logs", 0755)
//     if err != nil {
//         log.Fatal("Failed to create logs directory")
//     }

//     // Create or append to log file
//     currentTime := time.Now().Format("2006-01-02")
//     logFile, err := os.OpenFile(fmt.Sprintf("logs/sage_ai_%s.log", currentTime), 
//         os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//     if err != nil {
//         log.Fatal("Failed to open log file")
//     }

//     // Initialize loggers
//     InfoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
//     ErrorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
//     DebugLogger = log.New(logFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

//     // Also write to stdout
//     InfoLogger.SetOutput(os.Stdout)
//     ErrorLogger.SetOutput(os.Stdout)
//     DebugLogger.SetOutput(os.Stdout)
// }

// backend/go/pkg/logger/logger.go
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Define loggers
var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
)

// Initialize loggers at package level
func init() {
	// Create logs directory if it doesn't exist
	logsDir := "./logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		log.Fatalf("Failed to create logs directory: %v", err)
	}

	// Create log file with today's date
	today := time.Now().Format("2006-01-02")
	logFileName := fmt.Sprintf("sage_ai_%s.log", today)
	logFilePath := filepath.Join(logsDir, logFileName)

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Set up loggers with multi-writer to output to both file and console
	infoWriter := io.MultiWriter(os.Stdout, logFile)
	errorWriter := io.MultiWriter(os.Stderr, logFile)
	debugWriter := io.MultiWriter(os.Stdout, logFile)

	InfoLogger = log.New(infoWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(debugWriter, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Write startup message
	InfoLogger.Println("Logger initialized")
}