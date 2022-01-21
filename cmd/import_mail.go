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
	"errors"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
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
func importMail(f io.Reader, d datastore.Saver) (string, error) {
	var rawMail Mails

	// convert import file into a mails struct
	err := parseFromFile(f, &rawMail)
	if err != nil {
		return "", fmt.Errorf("failed to parse input file: %w", err)
	}

	if len(rawMail.Mails) == 0 {
		return "", errors.New("no mail entries in import file")
	}

	mails := UniqueMails(rawMail.Mails)
	importCount := len(mails)

	// conduct import as a transaction
	tx, err := d.Begin(true)
	if err != nil {
		return "", fmt.Errorf("error beginning datastore transaction: %w", err)
	}
	defer func() {
		if errRollback := tx.Rollback(); err != nil { // only log a rollback error if an error was encountered when saving
			log.Error("failed to rollback datastore transaction: ", errRollback)
		}
	}()

	// datastore needs to add one listing at a time, walk through imported listings and save one by one
	for _, r := range mails {
		mail := r

		err = tx.Save(&mail)
		if err != nil {
			log.Warn("failed to import record:", err)
			importCount--
		}

		log.WithFields(log.Fields{
			"listing": fmt.Sprintf("%+v", mail),
		}).Debug("imported record")
	}
	log.WithFields(log.Fields{
		"import_count": importCount,
		"read_count":   len(rawMail.Mails),
	}).Info("completed importing records")

	if errCommit := tx.Commit(); errCommit != nil {
		return "", fmt.Errorf("error committing records to datastore: %w", errCommit)
	}

	// Tell user how many records were imported.
	return fmt.Sprintf("Imported %d/%d mail records.", importCount, len(rawMail.Mails)), nil
}

// UniqueMails returns the passed in slice of mail with at most one of each mail. Mail order is
// preserved by first occurrence in initial slice.
func UniqueMails(rawMails []Mail) []Mail {
	keys := make(map[Mail]bool)
	cleanMails := []Mail{}

	for _, mail := range rawMails {
		if _, found := keys[mail]; !found {
			keys[mail] = true
			cleanMails = append(cleanMails, mail)
		}
	}

	return cleanMails
}
