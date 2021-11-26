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

package cmd

import (

	//nolint:gosec // not using this for security purposes
	"crypto/md5"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
)

// Mail contains relevant information for correspondence.
type Mail struct {
	ID       int `storm:"id,increment"`
	Ref      string
	Sender   int
	Receiver int
	Date     time.Time
	Link     string
}

var (
	sender   int
	receiver int
	date     string
	link     string
)

// NewMailCmd creates a mail command.
func NewMailCmd() *cobra.Command {
	// cmd represents the mail command
	cmd := &cobra.Command{
		Use:   "mail",
		Short: "Tracks letters sent to/from penpals",
		RunE:  RunMailCmd,
	}

	dm := viper.GetInt("member")
	cmd.Flags().IntVarP(&sender, "sender", "s", dm, "Correspondence sender.")
	cmd.Flags().IntVarP(&receiver, "receiver", "r", dm, "Correspondence receiver.")
	cmd.Flags().StringVarP(&date, "date", "d", time.Now().Format("2006-01-02"), "Correspondence date.")
	cmd.Flags().StringVarP(&link, "link", "l", "", "Link to listing ID or previous correspondence. 'L' prefix for listing entry, 'M' prefix for mail")

	return cmd
}

// RunMailCmd implements functionality of a mail command.
func RunMailCmd(cmd *cobra.Command, args []string) error {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
		}).Error(`Unable to determine time location.`)
		return errors.New("failed to add correspondence")
	}

	d, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
			"date":    date,
		}).Error("Invalid date argument.")
		return errors.New("failed to add correspondence")
	}

	// TODO: validate link is valid

	hSrc := fmt.Sprint(sender, receiver, d)
	ref := fmt.Sprintf("%x", md5.Sum([]byte(hSrc))) //nolint:gosec // not using this for security purposes

	m := Mail{
		Ref:      ref[26:], // grab the last 6 characters of the hash
		Sender:   sender,
		Receiver: receiver,
		Date:     d,
		Link:     link,
	}

	dsManager, err := datastore.New(viper.GetString("datastore.filename"))
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
		}).Error("Failed to open datastore.")
		return errors.New("failed to add correspondence")
	}
	defer dsManager.Stop()

	err = dsManager.Save(&m)
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
		}).Error("Failed to save correspondence.")
		return errors.New("failed to add correspondence")
	}

	log.WithFields(log.Fields{
		"command":  cmd.Name(),
		"ref":      m.Ref,
		"sender":   m.Sender,
		"receiver": m.Receiver,
		"date":     m.Date.Local().Format("2006-01-02"),
		"link":     m.Link,
	}).Info("Added mail entry.")

	cmd.Printf("Added mail. Reference: %s\n", m.Ref)
	return nil
}

func init() {
	// Search config in application directory with name ".ogma" (without extension).
	appFS := afero.NewOsFs()
	viper.SetFs(appFS)
	viper.AddConfigPath("./")
	viper.AddConfigPath("$HOME/")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".ogma")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(log.Fields{
			"configFile": "ogma",
			"err":        err,
		}).Warn("No config file found.")
	}

	rootCmd.AddCommand(NewMailCmd())
}
