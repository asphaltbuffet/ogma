// Application which greets you.
package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/cmd"
)

func main() {
	log.Trace("Started ogma.")
	defer log.Trace("Closed ogma.")

	appFS := afero.NewOsFs()
	InitConfig(appFS, ".ogma")

	cmd.Execute()
}

// InitConfig sets up the viper configuration.
func InitConfig(fs afero.Fs, cf string) {
	// Search config in application directory with name ".ogma" (without extension).
	viper.SetFs(fs)
	viper.AddConfigPath("./")
	viper.AddConfigPath("$HOME/")
	viper.SetConfigType("yaml")
	viper.SetConfigName(cf)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(log.Fields{
			"configFile": cf,
			"err":        err,
		}).Warn("No config file found.")
	}

	loggingLevel, err := log.ParseLevel(viper.GetString("logging.level"))
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Warn("Failed to read logging level from config.")
		loggingLevel = log.InfoLevel // default to info level logging
	}

	log.WithFields(log.Fields{"level": loggingLevel}).Debug("Setting log level.")
	log.SetLevel(loggingLevel)
}

func init() {
	// TODO: add a logging init here? 2021-11-19 BL
}
