package utils

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
)

func InitLogger() {
	logFilePath, _ := GetEnviromentVariable("DAEMON_LOG_FILE")

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil || logFilePath == "" {
		InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		WarnLogger = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
		ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		WarnLogger.Print("DAEMON_LOG_FILE environment variable not found. Logs are not stored in log file.")
	} else {
		InfoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
		WarnLogger = log.New(logFile, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
		ErrorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		InfoLogger.Printf("Logging to file: %s", logFilePath)
	}

}
