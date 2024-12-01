package logger

import (
	"log"
	"os"
)

var (
	LogInfo  *log.Logger
	LogDebug *log.Logger
	LogWarn  *log.Logger
	LogError *log.Logger
)

func init() {

	//output, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	//if err != nil {
	//	log.Fatal(err)
	//}

	output := os.Stdout

	LogInfo = log.New(output, "INFO:    ", log.Ldate|log.Ltime)
	LogDebug = log.New(output, "DEBUG:   ", log.Ldate|log.Ltime)
	LogWarn = log.New(output, "WARNING: ", log.Ldate|log.Ltime)
	LogError = log.New(output, "ERROR:   ", log.Ldate|log.Ltime|log.Lshortfile)
}
