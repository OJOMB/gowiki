package wikiLog

import (
	"gowiki/wikiUtils"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var ONCE sync.Once

type Logger struct {
	filename string
	Warning  *log.Logger
	Info     *log.Logger
	Error    *log.Logger
	Fatal    *log.Logger
}

func GetWikiLogger() *Logger {
	var globalLogger *Logger
	ONCE.Do(func() { globalLogger = initLogger() })
	return globalLogger
}

func initLogger() *Logger {
	filename := "log-" + strings.ReplaceAll(time.Now().UTC().Format(time.UnixDate), " ", "_")
	filepath := wikiUtils.ConstructPath(
		[]string{"wikiLog", "wikiLogFiles", filename},
		".log",
	)
	file, err := os.Create(filepath)
	if err != nil {
		// log.Fatal calls os.Exit(1) after writing to stderr
		log.Fatal(err)
	}

	var WarningLogger *log.Logger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	var InfoLogger *log.Logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	var ErrorLogger *log.Logger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	var FatalLogger *log.Logger = log.New(file, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{
		filename: filename,
		Warning:  WarningLogger,
		Info:     InfoLogger,
		Error:    ErrorLogger,
		Fatal:    FatalLogger,
	}
}
