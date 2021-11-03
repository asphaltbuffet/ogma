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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
)

// RunImportListings adds one to many listings to the datastore from a file.
func RunImportListings(c *cobra.Command) error {
	dsManager, err := datastore.New(viper.GetString("datastore.filename"))
	if err != nil {
		log.Error("Datastore manager failure.")
		return err
	}

	defer dsManager.Stop()
	if err != nil {
		log.Error("Failed to save to db: ")
		return err
	}

	newListing, err := ParseInputForListing(c)
	if err != nil {
		log.Error("Failed to save to db: ")
		return err
	}

	log.WithFields(log.Fields{
		"cmd": "listings.import",
	}).Info("Adding a listing.")

	err = dsManager.Store.Save(&newListing)
	if err != nil {
		log.Error("Failed to save new listings.")
		return err
	}

	c.Println("Added a listing.")
	c.Println(newListing.Render())

	return nil
}
