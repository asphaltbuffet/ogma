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
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
)

// importCmd represents the base command when called without any subcommands.
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Bulk import mailrecords.",
}

func initImportFile(f string) (io.ReadCloser, datastore.WriteCloser, error) {
	jsonFile, err := os.Open(filepath.Clean(f))
	if err != nil {
		log.Error("failed to open import file: ", err)

		return nil, nil, fmt.Errorf("failed to open import file: %w", err)
	}

	log.Debug("Successfully opened import file.")

	dsManager, err := datastore.New(viper.GetString("datastore.filename"))
	if err != nil {
		log.Error("failed to access datastore: ", err)
		return nil, nil, fmt.Errorf("failed to access datastore: %w", err)
	}

	return jsonFile, dsManager, nil
}

func init() {
	rootCmd.AddCommand(importCmd)
}
