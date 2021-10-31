// Application which greets you.
package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/cmd"
	"github.com/asphaltbuffet/ogma/pkg/datastore"
)

func main() {
	PrepConfigAndLogging()
	log.Info("Starting ogma...")
	cmd.Execute()

	dsManager, err := datastore.New("ogma.db")
	if err != nil {
		log.Fatal("Datastore manager failure.")
	}

	defer dsManager.Stop()

	// err = dsManager.Store.Save(&listing)
	if err != nil {
		log.Error("Failed to save to db: ")
	}
}

// PrepConfigAndLogging sets up configuration and logging for application.
func PrepConfigAndLogging() {
	// Search config in application directory with name ".ogma" (without extension).
	viper.AddConfigPath("./")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".ogma")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "No config file found:", viper.ConfigFileUsed())
	}

	// Create the log file if doesn't exist. And append to it if it already exists.
	// f, err := os.OpenFile(viper.GetString("logging.filename"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o600) //nolint:gomnd // setting file permissions
	// Formatter := new(log.TextFormatter)
	// You can change the Timestamp format. But you have to use the same date and time.
	// "2006-02-02 15:04:06" Works. If you change any digit, it won't work
	// ie "Mon Jan 2 15:04:05 MST 2006" is the reference time. You can't change it
	// Formatter.TimestampFormat = "2006-01-02 15:04:05"
	// Formatter.FullTimestamp = true
	// log.SetFormatter(Formatter)
	// if err != nil {
	// Cannot open log file. Logging to stderr
	//	fmt.Println(err)
	// } else {
	//	log.SetOutput(f)
	// }

	// TODO: Use Parse.LogLevel
	switch viper.GetString("logging.level") {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.WarnLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}
}
