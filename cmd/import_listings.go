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
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
	lstg "github.com/asphaltbuffet/ogma/pkg/listing"
)

const importListingCommandLongDesc = "Imports one-to-many LEX ads from a json file. This json file\n" +
	"should follow the format provided in the project 'examples' directory."

func init() {
	importCmd.AddCommand(NewImportListingCmd())
}

// NewImportListingCmd sets up an import subcommand.
func NewImportListingCmd() *cobra.Command {
	// cmd represents the import listing command.
	cmd := &cobra.Command{
		Use:     "listings [filename]",
		Short:   "Bulk import listing records.",
		Long:    importListingCommandLongDesc,
		Example: "ogma import listings listingImport.json",
		Run:     RunImportListingsCmd,
	}

	return cmd
}

// RunImportListingsCmd performs action associated with listings-import application command.
func RunImportListingsCmd(cmd *cobra.Command, args []string) {
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

	listOut, err := ImportListings(jsonFile, dsManager)
	if err != nil {
		log.Error("failed to import listing records: ", err)
		cmd.PrintErrln("failed to import listing records: ", err)
		return
	}

	cmd.Println(listOut)
}

// ImportListings adds one to many listings to the datastore from a file.
func ImportListings(f io.Reader, d datastore.Saver) (string, error) {
	// convert import file into a listings struct
	var rawListings lstg.Listings
	err := parseFromFile(f, &rawListings)
	if err != nil {
		return "", fmt.Errorf("failed to parse input file: %w", err)
	}

	listings := UniqueListings(rawListings.Listings)
	importCount := len(rawListings.Listings)

	// conduct import as a transaction
	tx, err := d.Begin(true)
	if err != nil {
		return "", fmt.Errorf("error beginning datastore transaction: %w", err)
	}
	defer func() {
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Error("failed to rollback datastore transaction: ", errRollback)
		}
	}()

	// datastore needs to add one listing at a time, walk through imported listings and save one by one
	for _, l := range listings {
		listing := l

		err = tx.Save(&listing)
		if err != nil {
			log.WithFields(log.Fields{
				"cmd":     "import",
				"listing": fmt.Sprintf("%+v", listing),
			}).Warn("failed to import record:", err)
			importCount--
		}

		log.WithFields(log.Fields{
			"cmd":     "import",
			"listing": fmt.Sprintf("%+v", listing),
		}).Debug("imported record")
	}
	log.WithFields(log.Fields{
		"cmd":          "import",
		"import_count": importCount,
		"read_count":   len(rawListings.Listings),
	}).Info("completed importing records")

	if errCommit := tx.Commit(); errCommit != nil {
		return "", fmt.Errorf("error committing records to datastore: %w", errCommit)
	}

	// Tell user how many records were imported.
	return fmt.Sprintf("Imported %d/%d listing records.", importCount, len(rawListings.Listings)), nil
}

// UniqueListings returns the passed in slice of listings with at most one of each listing. Listing order is
// preserved by first occurrence in initial slice.
func UniqueListings(rawListings []lstg.Listing) []lstg.Listing {
	log.WithFields(log.Fields{
		"cmd":   "import",
		"count": len(rawListings),
	}).Debug("deduplicating records")

	// Nothing to process. Return quickly.
	if len(rawListings) <= 1 {
		log.WithFields(log.Fields{
			"cmd":   "import",
			"count": 0,
		}).Debug("deduplication complete")
		return rawListings
	}

	keys := make(map[lstg.Listing]bool)
	cleanListings := []lstg.Listing{}

	for _, listing := range rawListings {
		if _, found := keys[listing]; !found {
			keys[listing] = true
			cleanListings = append(cleanListings, listing)
		}
	}

	log.WithFields(log.Fields{
		"cmd":   "import",
		"count": len(cleanListings),
	}).Debug("deduplication complete")

	return cleanListings
}
