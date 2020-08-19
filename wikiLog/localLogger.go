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

type LocalLogger struct {
	filename string
	warning  *log.Logger
	info     *log.Logger
	err      *log.Logger
	fatal    *log.Logger
}

func (l *LocalLogger) Info(msg ...interface{}) {
	passToPrint(l.info, msg...)
}

func (l *LocalLogger) Warn(msg ...interface{}) {
	passToPrint(l.warning, msg...)
}

func (l *LocalLogger) Error(msg ...interface{}) {
	passToPrint(l.err, msg...)
}

func (l *LocalLogger) Fatal(msg ...interface{}) {
	passToPrint(l.fatal, msg...)
}

func passToPrint(logger *log.Logger, args ...interface{}) {
	if len(args) > 1 {
		logger.Printf(args[0].(string)+"\n", args[1:]...)
		return
	}
	logger.Println(args[0])
}

func GetLocalLogger() *LocalLogger {
	var globalLocalLogger *LocalLogger
	ONCE.Do(func() { globalLocalLogger = initLocalLogger() })
	return globalLocalLogger
}

func initLocalLogger() *LocalLogger {
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

	return &LocalLogger{
		filename: filename,
		warning:  WarningLogger,
		info:     InfoLogger,
		err:      ErrorLogger,
		fatal:    FatalLogger,
	}
}
