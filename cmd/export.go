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
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
	lstg "github.com/asphaltbuffet/ogma/pkg/listing"
)

const exportCommandLongDesc = "The export command exports records from the datastore to json format. These files can be reimported."

var (
	exportType string
	exportFile string
)

// NewExportCmd creates an export command.
func NewExportCmd() *cobra.Command {
	// cmd represents the export command
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export stored records",
		Long:  exportCommandLongDesc,
		Run:   RunExportCmd,
	}

	cmd.Flags().StringVarP(&exportType, "record", "r", "all", "type of record to export ('mail', 'listing', or 'all')")
	cmd.Flags().StringVarP(&exportFile, "outfile", "o", "export.json", "file to export records to")

	return cmd
}

// RunExportCmd performs action associated with export command.
func RunExportCmd(cmd *cobra.Command, args []string) {
	switch exportType {
	case "all":
		if err := exportAll(); err != nil {
			cmd.Println("error exporting all data: ", err)
			log.Error("error exporting all data: ", err)
		}
	case "listing":
		if err := exportListing(); err != nil {
			cmd.Println("error exporting listing data: ", err)
			log.Error("error exporting listing data: ", err)
		}
	case "mail":
		if err := exportMail(); err != nil {
			cmd.Println("error exporting mail data: ", err)
			log.Error("error exporting mail data: ", err)
		}
	default:
		cmd.PrintErrln("invalid option: ", exportType)
		return
	}

	cmd.Println("successfully exported data")
}

func exportMail() error {
	ds, err := datastore.New(viper.GetString("datastore.filename"))
	if err != nil {
		return fmt.Errorf("error accessing datastore: %w", err)
	}

	var mailRecords Mails
	err = ds.All(&mailRecords.Mails)
	if err != nil {
		return fmt.Errorf("error getting mail records: %w", err)
	}

	// TODO: need to edit all IDs to be '0' so that reimporting doesn't
	// overwrite datastore values OR change import to have an overwrite
	// flag to provide edit functionality
	mailData, err := json.Marshal(mailRecords)
	if err != nil {
		return fmt.Errorf("error marshaling mail records: %w", err)
	}

	err = os.WriteFile(exportFile, mailData, 0o600)
	if err != nil {
		return fmt.Errorf("error writing mail data export: %w", err)
	}

	return nil
}

func exportListing() error {
	ds, err := datastore.New(viper.GetString("datastore.filename"))
	if err != nil {
		return fmt.Errorf("error accessing datastore: %w", err)
	}

	var listingRecords lstg.Listings
	err = ds.All(&listingRecords.Listings)
	if err != nil {
		return fmt.Errorf("error getting listing records: %w", err)
	}

	// TODO: need to edit all IDs to be '0' so that reimporting doesn't
	// overwrite datastore values OR change import to have an overwrite
	// flag to provide edit functionality
	listingData, err := json.Marshal(listingRecords)
	if err != nil {
		return fmt.Errorf("error marshaling listing records: %w", err)
	}

	// TODO: write to file instead of stdout
	fmt.Println(string(listingData))

	return nil
}

func exportAll() error {
	if err := exportListing(); err != nil {
		return fmt.Errorf("error with exporting listings when exporting all records: %w", err)
	}

	if err := exportMail(); err != nil {
		return fmt.Errorf("error with exporting mail when exporting all records: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(NewExportCmd())
}
