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
	"io"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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

const (
	// MaxHashLength is the maximum hash length for md5 checksum.
	MaxHashLength = 32

	// MinHashLength is the minimum hash length for md5 checksum.
	MinHashLength = 0

	// RefLength is the default reference length.
	RefLength = 6
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
	cmd.Flags().IntP("sender", "s", dm, "Correspondence sender.")
	cmd.Flags().IntP("receiver", "r", dm, "Correspondence receiver.")
	cmd.Flags().StringP("date", "d", time.Now().Format("2006-01-02"), "Correspondence date.")
	cmd.Flags().StringP("link", "l", "", "Link to listing ID or previous correspondence. 'L' prefix for listing entry, 'M' prefix for mail")
	cmd.Flags().IntP("length", "L", RefLength, "Correspondence receiver.")

	return cmd
}

// RunMailCmd implements functionality of a mail command.
func RunMailCmd(cmd *cobra.Command, args []string) error {
	m, err := MailFromArgs(cmd)
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
			"args":    args,
			"err":     err,
		}).Error("failed to parse arguments")
		return err
	}

	dsManager, err := datastore.New(viper.GetString("datastore.filename"))
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
		}).Error("failed to open datastore")
		return errors.New("datastore: failed to add correspondence")
	}
	defer dsManager.Stop()

	err = dsManager.Save(&m)
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
		}).Error("failed to save correspondence")
		return errors.New("save: failed to add correspondence")
	}

	log.WithFields(log.Fields{
		"command":  cmd.Name(),
		"ref":      m.Ref,
		"sender":   m.Sender,
		"receiver": m.Receiver,
		"date":     m.Date.Local().Format("2006-01-02"),
		"link":     m.Link,
	}).Info("added mail entry")

	cmd.Printf("Added mail. Reference: %s\n", m.Ref)
	return nil
}

// MailFromArgs creates a new mail object from command arguments.
func MailFromArgs(cmd *cobra.Command) (Mail, error) {
	s, err := cmd.Flags().GetInt("sender")
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
			"args":    cmd.Args,
		}).Warn("failed to get sender argument")
	}

	r, err := cmd.Flags().GetInt("receiver")
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
			"args":    cmd.Args,
		}).Warn("failed to get receiver argument")
	}

	m := Mail{
		Sender:   s,
		Receiver: r,
	}

	date, err := cmd.Flags().GetString("date")
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
			"args":    cmd.Args,
		}).Warn("failed to get date argument")
	}

	d, err := ValidateDate(date)
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
			"date":    date,
		}).Error("failed to validate date")
		return Mail{}, errors.New("date: failed to add correspondence")
	}
	m.Date = d

	link, err := cmd.Flags().GetString("link")
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
			"args":    cmd.Args,
		}).Warn("failed to get link argument")
	}
	// TODO: validate link is valid
	m.Link = link

	m.Ref = MailHash(m, RefLength)

	return m, nil
}

// MailHash creates a 'unique' hash of the sender, receiver, and mail date.
func MailHash(m Mail, l int) string {
	if l > MaxHashLength {
		l = MaxHashLength
	} else if l < MinHashLength {
		l = MinHashLength
	}

	h := md5.New() //nolint:gosec // not using this for security purposes
	padding := "qwertyuiopasdfghjklzxcvbnm1234567890"
	hSrc := fmt.Sprint(m.Sender, m.Receiver, m.Date.Local().Format("2006-01-02"))
	if _, err := io.WriteString(h, padding); err != nil {
		log.WithFields(log.Fields{
			"pre-hash": hSrc,
		}).Warn(`failed to calculate reference`)
		return ""
	}

	ref := fmt.Sprintf("%x", md5.Sum([]byte(hSrc))) //nolint:gosec // not using this for security purposes
	log.WithFields(log.Fields{
		"pre-hash":  hSrc,
		"full-hash": ref,
	}).Debug(`calculated reference hash`)

	return ref[len(ref)-l:]
}

// ValidateDate checks date string format and parses with local time location.
func ValidateDate(dd string) (time.Time, error) {
	// hardcoded location for now. could be a configuration later.
	location := "Local"
	loc, err := time.LoadLocation(location)
	if err != nil {
		log.WithFields(log.Fields{
			"time_location": location,
		}).Warn(`failed to load time location`)

		// set to UTC if cannot get local location
		loc = time.UTC
	}

	d, err := time.ParseInLocation("2006-01-02", dd, loc)
	if err != nil {
		log.WithFields(log.Fields{
			"date":          dd,
			"time_location": location,
		}).Error("invalid date argument")
		return time.Time{}, err
	}

	return d, nil
}

// RenderMail returns a pretty formatted listing as table.
func RenderMail(mm []Mail, s ...string) string {
	// empty string if there are no listings to render
	if len(mm) == 0 {
		return ""
	}

	mt := table.NewWriter()

	mt.SetTitle("Correspondence Matches:")

	mt.AppendHeader(table.Row{
		"Reference",
		"Sender",
		"Receiver",
		"Date",
		"Link",
	})

	for _, m := range mm {
		mt.AppendRow([]interface{}{
			m.Ref,
			m.Sender,
			m.Receiver,
			m.Date.Local().Format("2006-01-02"),
			m.Link,
		})
	}

	columnConfigs := []table.ColumnConfig{
		{
			Name:  "Sender",
			Align: text.AlignRight,
		},
		{
			Name:  "Receiver",
			Align: text.AlignRight,
		},
		{
			Name:  "Ref",
			Align: text.AlignCenter,
		},
		{
			Name:  "Date",
			Align: text.AlignRight,
		},
		{
			Name:  "Link",
			Align: text.AlignCenter,
		},
	}
	mt.SetColumnConfigs(columnConfigs)

	mt.SortBy([]table.SortBy{
		{Name: "Date", Mode: table.Asc},
		{Name: "Ref", Mode: table.Asc},
	})

	style, _ := GetStyle(s)
	mt.SetStyle(style)

	return mt.Render()
}

// GetStyle converts a []string into a valid rendering style.
func GetStyle(s []string) (table.Style, error) {
	// set default style as bright
	sa := "bright"

	// if too many args, just use the first one
	if len(s) > 0 {
		if len(s) > 1 {
			log.WithFields(log.Fields{
				"args": s,
			}).Warn("Too many styles passed in. Using first argument only.")
		}

		sa = s[0]
	}

	// check style
	switch sa {
	case "bright":
		return table.StyleColoredBright, nil
	case "light":
		return table.StyleLight, nil
	case "default":
		return table.StyleDefault, nil
	default:
		log.WithFields(log.Fields{
			"style": sa,
		}).Error("Invalid table style.")
		return table.StyleDefault, errors.New("invalid argument")
	}
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
		}).Warn("no config file found")
	}

	rootCmd.AddCommand(NewMailCmd())
}
