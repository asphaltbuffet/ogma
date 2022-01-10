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
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const deleteCommandLongDesc = "The delete command removes the datastore file. This action cannot be undone!"

// NewDeleteCmd creates a delete command.
func NewDeleteCmd() *cobra.Command {
	// cmd represents the delete command
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete records",
		Long:  deleteCommandLongDesc,
		Run:   RunDeleteCmd,
	}

	cmd.Flags().BoolP("all", "a", false, "remove all entries")

	return cmd
}

func init() {
	rootCmd.AddCommand(NewDeleteCmd())
}

// RunDeleteCmd performs action associated with delete command.
func RunDeleteCmd(cmd *cobra.Command, args []string) {
	isAll, err := cmd.Flags().GetBool("all")
	if err != nil {
		cmd.Println("error getting 'all' flag: ", err)
	}

	if isAll {
		if err := clearData(); err != nil {
			cmd.Println("error deleting all data: ", err)
			log.Error("error deleting all data: ", err)
		}

		cmd.Println("all data records have been removed")
	}
}

func clearData() error {
	dsFile := viper.GetString("datastore.filename")

	if _, err := os.Stat(dsFile); err != nil {
		return fmt.Errorf("error finding datastore file: %w", err)
	}

	log.Infof("clearing data from datastore filename=%s", dsFile)

	if err := os.Remove(dsFile); err != nil {
		return fmt.Errorf("error deleting file: %w", err)
	}

	return nil
}
