//go:build !windows

package main

import (
	"log/syslog"
	"runtime"

	log "github.com/sirupsen/logrus"
	lSyslog "github.com/sirupsen/logrus/hooks/syslog"
)

func init() {
	// syslog hook doesn't work in windows.
	if runtime.GOOS != "windows" {
		hook, err := lSyslog.NewSyslogHook("", "", syslog.LOG_INFO, "")
		if err != nil {
			log.Warn("Unable to set up syslog hook.")
		}
		log.AddHook(hook)
	}
}
