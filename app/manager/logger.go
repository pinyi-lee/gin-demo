package manager

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	loggerManager *LoggerManager
)

func GetLogger() *LoggerManager {
	return loggerManager
}

type LoggerManager struct {
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger

	config LoggerConfig
}

type LoggerConfig struct {
	LogLevel string
}

func (manager *LoggerManager) Setup(config LoggerConfig) (err error) {

	logLevel := strings.ToUpper(config.LogLevel)

	debugHandle := ioutil.Discard
	infoHandle := ioutil.Discard
	warnHandle := ioutil.Discard
	errorHandle := os.Stderr

	if logLevel == "DEBUG" {
		debugHandle = os.Stdout
		infoHandle = os.Stdout
		warnHandle = os.Stdout
	}
	if logLevel == "INFO" {
		infoHandle = os.Stdout
		warnHandle = os.Stdout
	}
	if logLevel == "WARN" {
		warnHandle = os.Stdout
	}

	debugLog := log.New(debugHandle, "\u001b[46;1m\u001b[37mDEB:\u001b[0m ", log.Ldate|log.Lmicroseconds|log.LUTC|log.Llongfile)
	infoLog := log.New(infoHandle, "\u001b[42;1m\u001b[37mINF:\u001b[0m ", log.Ldate|log.Lmicroseconds|log.LUTC|log.Llongfile)
	warnLog := log.New(warnHandle, "\u001b[43;1m\u001b[38mWARN:\u001b[0m ", log.Ldate|log.Lmicroseconds|log.LUTC|log.Llongfile)
	errorLog := log.New(errorHandle, "\u001b[41;1m\u001b[37mERR:\u001b[0m ", log.Ldate|log.Lmicroseconds|log.LUTC|log.Llongfile)

	loggerManager = &LoggerManager{
		Debug:  debugLog,
		Info:   infoLog,
		Warn:   warnLog,
		Error:  errorLog,
		config: config,
	}

	return nil
}

func (manager *LoggerManager) Close() {

}
