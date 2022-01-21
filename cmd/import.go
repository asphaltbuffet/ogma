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
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
)

// importCmd represents the base command when called without any subcommands.
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Bulk import records.",
	Run:   RunImportCmd,
}

func init() {
	rootCmd.AddCommand(importCmd)
}

// RunImportCmd performs action associated with import application command.
func RunImportCmd(cmd *cobra.Command, args []string) {
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
		log.Error("error initializing import: ", err)
		cmd.PrintErrln("error initializing import: ", err)
		return
	}

	importSummary, err := importListings(jsonFile, dsManager)
	if err != nil {
		log.Error("failed to import listing records: ", err)
		cmd.PrintErrln("failed to import listing records: ", err)
		return
	}

	cmd.Println(importSummary)

	importSummary, err = importMail(jsonFile, dsManager)
	if err != nil {
		log.Error("failed to import mail records: ", err)
		cmd.PrintErrln("failed to import mail records: ", err)
		return
	}

	cmd.Println(importSummary)
}

// initImportFile is shared initialization for all import types and datastore.
func initImportFile(f string) (io.ReadCloser, datastore.SaveStopper, error) {
	jsonFile, err := os.Open(filepath.Clean(f))
	if err != nil {
		log.Error("failed to open import file: ", err)

		return nil, nil, fmt.Errorf("failed to open import file: %w", err)
	}

	log.Debug("Successfully opened import file.")

	dsManager, err := datastore.New(viper.GetString("datastore.filename"))
	if err != nil {
		if closeErr := jsonFile.Close(); closeErr != nil {
			log.Error("failed to close import file: ", closeErr)
		}

		log.Error("failed to access datastore: ", err)

		return nil, nil, fmt.Errorf("failed to access datastore: %w", err)
	}

	return jsonFile, dsManager, nil
}

// parseFromFile unmarshalls json into a Listings struct.
func parseFromFile(j io.Reader, value interface{}) error {
	if j == nil {
		return errors.New("argument cannot be nil")
	}

	// read our opened jsonFile as a byte array.
	byteValue, err := afero.ReadAll(j)
	if err != nil {
		return fmt.Errorf("failed to read import file: %w", err)
	}

	err = json.Unmarshal(byteValue, value)
	if err != nil {
		return fmt.Errorf("failed to unmarshall import file: %w", err)
	}

	return nil
}
