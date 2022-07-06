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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
)

// Mails is a container for multiple mail objects.
type Mails struct {
	Mails []Mail `json:"mails"`
}

// Mail contains relevant information for correspondence.
type Mail struct {
	ID       int    `storm:"id,increment"`
	Ref      string `json:"reference"`
	Sender   int    `json:"sender"`
	Receiver int    `json:"receiver"`
	Date     string `json:"date"`
	Link     string `json:"link"`
}

const (
	// MaxHashLength is the maximum hash length for md5 checksum.
	MaxHashLength = 32

	// MinHashLength is the minimum hash length for md5 checksum.
	MinHashLength = 0

	// RefLength is the default reference length.
	RefLength = 6
)

const mailCommandLongDesc = "The mail command supports entering correspondence details and getting a reference number\n" +
	"back that can be used for tracking physical artifacts.\n\n" +
	"'Date' must be in the 'yyyy-mm-dd' format.\n" +
	"'Link' is optional. It must start with 'L' to link with an ad (using ID field from ad output) or 'M' to link to a correspondence reference."

var mailColumnConfigs = []table.ColumnConfig{
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

func init() {
	rootCmd.AddCommand(NewMailCmd())
}

// NewMailCmd creates a mail command.
func NewMailCmd() *cobra.Command {
	// cmd represents the mail command
	cmd := &cobra.Command{
		Use:     "mail",
		Short:   "Tracks letters sent to/from penpals",
		Long:    mailCommandLongDesc,
		Example: `ogma mail --sender=1234 --receiver=5678 --date=2021-11-02`,
		Run:     RunMailCmd,
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
func RunMailCmd(cmd *cobra.Command, args []string) {
	log.Debug("Running mail command. Config is: ", viper.ConfigFileUsed())
	m, err := mailFromArgs(cmd)
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
			"args":    args,
		}).Error("failed to parse arguments: ", err)
		cmd.PrintErrln("invalid input: ", err)
		return
	}

	dsManager, err := datastore.New(viper.GetString("datastore.filename"))
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
		}).Error("failed to open datastore: ", err)
		cmd.PrintErrln("Failed to access datastore: ", err)
		return
	}
	defer dsManager.Stop()

	err = dsManager.Save(&m)
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
		}).Error("unable to save mail: ", err)
		cmd.PrintErrln("Failed to save entry: ", err)
		return
	}

	log.WithFields(log.Fields{
		"command":  cmd.Name(),
		"ref":      m.Ref,
		"sender":   m.Sender,
		"receiver": m.Receiver,
		"date":     m.Date,
		"link":     m.Link,
	}).Info("added mail entry")

	cmd.Printf("Added mail. Reference: %s\n", m.Ref)
}

// mailFromArgs creates a new mail object from command arguments.
func mailFromArgs(cmd *cobra.Command) (Mail, error) {
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

	m.Date, err = ValidateDate(date)
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
			"date":    date,
		}).Error("failed to validate date")
		return Mail{}, errors.New("date: failed to add correspondence")
	}

	m.Link, err = cmd.Flags().GetString("link")
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
			"args":    cmd.Args,
		}).Warn("failed to get link argument")
	}

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
	hSrc := fmt.Sprint(m.Sender, m.Receiver, m.Date)
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
func ValidateDate(dd string) (string, error) {
	// hardcoded location for now. could be a configuration later.
	location := "Local"
	loc, err := time.LoadLocation(location)
	if err != nil {
		log.WithFields(log.Fields{
			"time_location": location,
		}).Warnf("unable to load local time location (using UTC): %v", err)

		// set to UTC if cannot get local location
		loc = time.UTC
	}

	d, err := time.ParseInLocation("2006-01-02", dd, loc)
	if err != nil {
		log.WithFields(log.Fields{
			"date":          dd,
			"time_location": location,
		}).Error("invalid date argument")
		return "", fmt.Errorf("date format must be 'yyyy-mm-dd': %w", err)
	}

	return d.Local().Format("2006-01-02"), nil
}

// RenderMail returns a pretty formatted listing as table.
func RenderMail(mm []Mail, p bool) string {
	// empty string if there are no listings to render
	if len(mm) == 0 {
		return "No correspondences found."
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
			m.Date,
			m.Link,
		})
	}

	mt.SetColumnConfigs(mailColumnConfigs)

	mt.SortBy([]table.SortBy{
		{Name: "Date", Mode: table.Asc},
		{Name: "Ref", Mode: table.Asc},
	})

	if p {
		mt.SetStyle(table.StyleColoredBright)
	}
	log.WithFields(log.Fields{
		"is_pretty": p,
	}).Debug("set rendering style")

	return mt.Render()
}
