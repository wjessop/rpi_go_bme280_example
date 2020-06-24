package main

import (
	logging "github.com/op/go-logging"

	"os"
)

var (
	log       = logging.MustGetLogger("climate")
	logFormat = logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.5s} %{id:03x}%{color:reset} %{message}`)
)

func createLogger(isDebug bool) {
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logFormatter := logging.NewBackendFormatter(logBackend, logFormat)
	logLeveled := logging.AddModuleLevel(logFormatter)
	logging.SetBackend(logLeveled)

	if isDebug {
		logLeveled.SetLevel(logging.DEBUG, "")
		log.Debugf("Setting log level to DEBUG")
	} else {
		logLeveled.SetLevel(logging.INFO, "")
	}
}
