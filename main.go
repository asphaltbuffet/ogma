// Application which greets you.
package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/cmd"
)

func main() {
	log.Trace("Started ogma.")
	defer log.Trace("Closed ogma.")

	cmd.Execute()
}

// InitConfig sets up the viper configuration.
func InitConfig(cf string) {
	// Search config in application directory with name ".ogma" (without extension).
	viper.AddConfigPath("./")
	viper.SetConfigType("yaml")
	viper.SetConfigName(cf)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "No config file found:", viper.ConfigFileUsed())
	}

	loggingLevel, err := log.ParseLevel(viper.GetString("logging.level"))
	if err != nil {
		loggingLevel = log.InfoLevel // default to info level logging
	}
	log.SetLevel(loggingLevel)
}

func init() {
	// TODO: add a logging init here? 2021-11-19 BL
	InitConfig(".ogma")
}
