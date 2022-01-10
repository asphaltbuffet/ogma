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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
)

const importMailCommandLongDesc = "Imports one-to-many correspondence records from a json file. This json\n" +
	"file should follow the format provided in the project 'examples' directory.\n\n" +
	"The reference field should be unique for each record. Manually creating this may increase the chance of duplicate values."

func init() {
	importCmd.AddCommand(NewImportMailCmd())
}

// NewImportMailCmd sets up an import subcommand.
func NewImportMailCmd() *cobra.Command {
	// cmd represents the import command.
	cmd := &cobra.Command{
		Use:     "mail [filename]",
		Short:   "Bulk import mail records.",
		Long:    importMailCommandLongDesc,
		Example: "ogma import mail mImport.json",
		Run:     RunImportMailCmd,
	}

	return cmd
}

// RunImportMailCmd performs action associated with mail-import application command.
func RunImportMailCmd(cmd *cobra.Command, args []string) {
	jsonFile, dsManager, err := initImportFile(args[0])
	// defer closing the import file until after we're done with it
	defer func() {
		if dsManager != nil {
			dsManager.Stop()
		}

		if jsonFile != nil {
			if closeErr := jsonFile.Close(); closeErr != nil {
				log.Error("failed to close import file: ", closeErr)
			}
		}
	}()
	if err != nil {
		log.Error("error initializing listings import: ", err)
		cmd.PrintErrln("error initializing listings import: ", err)
		return
	}

	mailOut, err := importMail(jsonFile, dsManager)
	if err != nil {
		log.Error("failed to import mail records: ", err)
		cmd.PrintErr("failed to import mail records: ", err)
		return
	}

	cmd.Println(mailOut)
}

// importMail adds one to many mail to the datastore from a file.
func importMail(f io.Reader, d datastore.Writer) (string, error) {
	// convert import file into a mails struct
	rawMail, err := parseMail(f)
	if err != nil {
		log.WithFields(log.Fields{"cmd": "import"}).Error("failed to parse input file: ", err)
		return "", fmt.Errorf("failed to parse input file: %w", err)
	}

	if len(rawMail) == 0 {
		log.Debug("no mail entries found to import")
		return "", errors.New("no mail entries in import file")
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

// parseMail unmarshalls json into a Mails struct.
func parseMail(j io.Reader) ([]Mail, error) {
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
