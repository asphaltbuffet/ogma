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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
)

// NewImportMailCmd sets up an import subcommand.
func NewImportMailCmd() *cobra.Command {
	// cmd represents the import command.
	cmd := &cobra.Command{
		Use:     "import",
		Short:   "Bulk import records.",
		Example: "ogma import somefile.json",
		Run:     RunImportCmd,
	}

	return cmd
}

// RunImportCmd performs action associated with listings-import application command.
func RunImportCmd(cmd *cobra.Command, args []string) {
	jsonFile, err := os.Open(args[0])
	if err != nil {
		log.Errorf("failed to open import file: %v", err)
		cmd.PrintErrf("failed to open import file: %v", err)
		return
	}

	log.Debug("Successfully opened import file.")

	// defer closing the import file until after we're done with it
	defer func() {
		err = jsonFile.Close()
		if err != nil {
			cmd.PrintErrf("failed to close import file: %v", err)
			log.Errorf("failed to close import file: %v", err)
			return
		}
	}()

	dsManager, err := datastore.New(viper.GetString("datastore.filename"))
	if err != nil {
		cmd.PrintErrf("failed to access datastore: %v", err)
		log.Errorf("failed to access datastore: %v", err)
		return
	}
	defer dsManager.Stop()

	mailOut, err := ImportMail(jsonFile, dsManager)
	if err != nil {
		log.Errorf("failed to import mail records: %v", err)
	}

	if mailOut != "" {
		cmd.Println(mailOut)
	}
}

// ImportMail adds one to many mail to the datastore from a file.
func ImportMail(f io.Reader, d datastore.Writer) (string, error) {
	// convert import file into a mails struct
	rawMail, err := ParseMail(f)
	if err != nil {
		log.WithFields(log.Fields{"cmd": "import"}).Error("failed to parse input file: ", err)
		return "", fmt.Errorf("failed to parse input file: %w", err)
	}

	if len(rawMail) == 0 {
		log.Debug("no mail entries found to import")
		return "", nil
	}

	mails := UniqueMails(rawMail)
	importCount := len(mails)

	// datastore needs to add one listing at a time, walk through imported listings and save one by one
	for _, r := range mails {
		mail := r

		err = d.Save(&mail)
		if err != nil {
			log.WithFields(log.Fields{
				"cmd": "import",
				// "mail": fmt.Sprintf("%v", mail),
			}).Warn("failed to import record:", err)
			importCount--
		}

		log.WithFields(log.Fields{
			"cmd":     "import",
			"listing": fmt.Sprintf("%+v", mail),
		}).Debug("imported record")
	}
	log.WithFields(log.Fields{
		"cmd":          "import",
		"import_count": importCount,
		"read_count":   len(rawMail),
	}).Info("completed importing records")

	// In all cases, tell user how many records were imported.
	return fmt.Sprintf("Imported %d/%d mail records.", importCount, len(rawMail)), nil
}

// ParseMail unmarshalls json into a Mails struct.
func ParseMail(j io.Reader) ([]Mail, error) {
	if j == nil {
		return []Mail{}, errors.New("argument cannot be nil")
	}

	// read our opened jsonFile as a byte array.
	byteValue, err := afero.ReadAll(j)
	if err != nil {
		log.WithFields(log.Fields{
			"cmd": "import",
		}).Error("failed to read import file:", err)
		return []Mail{}, fmt.Errorf("failed to read import file: %w", err)
	}

	var newMails Mails

	err = json.Unmarshal(byteValue, &newMails)
	if err != nil {
		log.WithFields(log.Fields{
			"cmd": "import",
		}).Error("failed to unmarshall import file:", err)
		return []Mail{}, fmt.Errorf("failed to unmarshall import file: %w", err)
	}

	return newMails.Mails, nil
}

// UniqueMails returns the passed in slice of mail with at most one of each mail. Mail order is
// preserved by first occurrence in initial slice.
func UniqueMails(rawMails []Mail) []Mail {
	log.WithFields(log.Fields{
		"cmd":   "import",
		"count": len(rawMails),
	}).Debug("deduplicating records")

	// Nothing to process. Return quickly.
	if len(rawMails) <= 1 {
		log.WithFields(log.Fields{
			"cmd":   "import",
			"count": 0,
		}).Debug("deduplication complete")
		return rawMails
	}

	keys := make(map[Mail]bool)
	cleanMails := []Mail{}

	for _, mail := range rawMails {
		if _, found := keys[mail]; !found {
			keys[mail] = true
			cleanMails = append(cleanMails, mail)
		}
	}

	log.WithFields(log.Fields{
		"cmd":   "import",
		"count": len(cleanMails),
	}).Debug("deduplication complete")

	return cleanMails
}

func init() {
	importCmd.AddCommand(NewImportMailCmd())
}
