/*
Copyright © 2021 Ben Lechlitner <otherland@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

// Package cmd contains all CLI commands used by the application.
package cmd

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// default const values for application.
const (
	DefaultMaxSearchResults  = 10
	DefaultMemberNumber      = 13401
	DefaultConfigFilename    = ".ogma"
	DefaultLoggingLevel      = "info"
	DefaultDatastoreFilename = "ogma.db"
)

var appFS afero.Fs

const rootCommandLongDesc = "Ogma is a tracking application for penpals using LEX magazine.\n" +
	"It stores a digital record of LEX magazine ads and allows the user to track letters\n" +
	"sent and received. Correspondence may also be linked to previous letters or LEX ads."

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:               "ogma",
	Version:           "1.1.1",
	Short:             "Ogma is a pen-pal tracking application.",
	Long:              rootCommandLongDesc,
	Args:              cobra.NoArgs,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	RunE: func(cmd *cobra.Command, args []string) error {
		// providing a run for root so that the PersistentPreRun will kick off to initialize config.
		return errors.New("no arguments provided")
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		InitConfig(appFS, DefaultConfigFilename)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// InitConfig sets up Viper and Logging.
func InitConfig(fs afero.Fs, cfg string) {
	log.Trace("initializing configuration and logging")

	appFS = afero.NewOsFs()

	viper.SetFs(appFS)

	// Search config in application directory with name ".ogma" (without extension).
	viper.AddConfigPath("./")
	viper.AddConfigPath("$HOME/")
	viper.SetConfigType("yaml")
	viper.SetConfigName(cfg)

	viper.SetEnvPrefix("OGMA")
	viper.AutomaticEnv() // read in environment variables that match

	viper.SetDefault("logging.level", DefaultLoggingLevel)
	viper.SetDefault("datastore.filename", DefaultDatastoreFilename)
	viper.SetDefault("search:max_results", DefaultMaxSearchResults)
	viper.SetDefault("member", DefaultMemberNumber)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Warn("unable to read config file: ", err)
	}

	loggingLevel, err := log.ParseLevel(viper.GetString("logging.level"))
	if err != nil {
		log.Warn("error parsing logging level: ", err)
	}

	log.SetLevel(loggingLevel)
	log.WithFields(log.Fields{"level": loggingLevel}).Debug("set log level")
}

func init() {
	// rootCmd.PersistentFlags().String("config", ".ogma", "Configuration file to use for application.")
}
