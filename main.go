// Application which greets you.
package main

import (
	"fmt"
	"os"

	"github.com/asdine/storm/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/cmd"
)

// A Listing contains relevant information for LEX listings.
type Listing struct {
	ID                  int `storm:"id,increment"`
	IssueNumber         int
	PageNumber          int
	IndexedMemberNumber int    `storm:"index"`
	IndexedCategory     string `storm:"index"`
	ListingText         string
}

func main() {
	PrepConfigAndLogging()
	log.Info("Starting ogma...")
	cmd.Execute()

	db, err := storm.Open("ogma.db")
	if err != nil {
		fmt.Println("Failed to open db: ", err)
	}

	defer func() {
		err = db.Close()
	}() // use function closure to allow checking error from deferred db.Close
	if err != nil {
		fmt.Println("Failed to close db: ", err)
	}

	listing := Listing{
		ID:                  1,
		IssueNumber:         56, //nolint:gomnd // preliminary dev magic number use
		PageNumber:          1,
		IndexedMemberNumber: 2989, //nolint:gomnd // preliminary dev magic number use
		IndexedCategory:     "Art & Photography",
		ListingText:         "Fingerpainting exchange.",
	}

	err = db.Save(&listing)
	if err != nil {
		fmt.Println("Failed to save to db: ", err)
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
