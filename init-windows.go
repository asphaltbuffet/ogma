//go:build windows

package main

import (
	log "github.com/sirupsen/logrus"
)

func init() {
	log.Debug("syslog logging not available in windows")
}
