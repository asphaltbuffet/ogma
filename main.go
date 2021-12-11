// Application which greets you.
package main

import (
	"log/syslog"
	"runtime"

	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"

	"github.com/asphaltbuffet/ogma/cmd"
)

func main() {
	cmd.Execute()
}

func init() {
	if runtime.GOOS != "windows" {
		hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
		if err != nil {
			log.Warn("Unable to set up syslog hook.")
		}
		log.AddHook(hook)
	}
}
