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
	"io"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
	lstg "github.com/asphaltbuffet/ogma/pkg/listing"
)

// importCmd represents the import command.
var importCmd = &cobra.Command{
	Use:     "import",
	Short:   "Bulk import records.",
	Long:    ``,
	Args:    cobra.ExactArgs(1),
	Example: "ogma import somefile.json -v",
	RunE:    RunImportCmd,
}

var verbose bool

// RunImportCmd performs action associated with listings-import application command.
func RunImportCmd(c *cobra.Command, args []string) error {
	jsonFile, err := os.Open(args[0])
	if err != nil {
		log.Error("Failed to open import file.")
		return err
	}

	log.Info("Successfully opened import file.")

	// defer closing the import file until after we're done with it
	defer func() {
		err = jsonFile.Close()
		if err != nil {
			log.Error("Failed to close import file.")
		}
	}()

	dsManager, err := datastore.New(viper.GetString("datastore.filename"))
	if err != nil {
		log.Error("Datastore manager failure.")
		return err
	}
	defer dsManager.Stop()

	out, err := Import(jsonFile, dsManager)

	if err == nil {
		if verbose {
			c.Println(out)
		}
	}

	return err
}

// Import adds one to many listings to the datastore from a file.
func Import(f io.Reader, d datastore.Writer) (string, error) {
	// verify that the import reader is valid
	if f == nil {
		return "", errors.New("import : import reader cannot be nil")
	}

	// convert import file into a listings struct
	rawListings, err := ParseImportInput(f)
	if err != nil {
		log.WithFields(log.Fields{"cmd": "listings.import"}).Error("Failed to import listings.")
		return "", err
	}

	listings := UniqueListings(rawListings)

	// datastore needs to add one listing at a time, walk through imported listings and save one by one
	for _, l := range listings {
		listing := l

		err = d.Save(&listing)
		if err != nil {
			log.WithFields(log.Fields{"cmd": "listings.import", "count": len(listings)}).Error("Failed datastore save.")
			return "", err
		}
		log.WithFields(log.Fields{"cmd": "listings.import", "listing": listing}).Debug("Imported listing.")
	}
	log.WithFields(log.Fields{"cmd": "listings.import", "count": len(listings)}).Info("Imported listings.")

	// In all cases, tell user how many records were imported.
	output := "Imported " + strconv.Itoa(len(listings)) + " record"
	if len(listings) == 1 {
		output += ".\n"
	} else {
		output += "s.\n"
	}

	return output + lstg.Render(listings), nil
}

// ParseImportInput unmarshalls json into a Listings struct.
func ParseImportInput(j io.Reader) ([]lstg.Listing, error) {
	// verify that the parameter is valid
	if j == nil {
		return []lstg.Listing{}, errors.New("import : parameter cannot be nil")
	}

	// read our opened jsonFile as a byte array.
	byteValue, err := afero.ReadAll(j)
	if err != nil {
		log.WithFields(log.Fields{
			"cmd": "listings.import",
		}).Error("Failed to unmarshall import file.")
		return []lstg.Listing{}, err
	}

	var ll lstg.Listings

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &ll)
	if err != nil {
		log.WithFields(log.Fields{
			"cmd": "listings.import",
		}).Error("Failed to unmarshall import file.")
		return []lstg.Listing{}, err
	}

	return ll.Listings, nil
}

// AddListing adds a single listing to the datastore.
func AddListing(rawRecords []lstg.Listing, ds datastore.Writer) (string, error) {
	// if there's nothing to add, return quickly
	if len(rawRecords) == 0 {
		return "", nil
	}

	cleanRecords := UniqueListings(rawRecords)

	log.WithFields(log.Fields{
		"cmd":        "listings.add",
		"count":      len(cleanRecords),
		"duplicates": len(rawRecords) - len(cleanRecords),
	}).Info("Adding listing(s).")

	for _, record := range cleanRecords {
		// copy loop variable so i can accurately reference it for saving
		listing := record
		err := ds.Save(&listing)
		if err != nil {
			log.Error("Failed to save new listing.")
			return "", err
		}
	}

	return lstg.Render(cleanRecords), nil
}

// UniqueListings returns the passed in slice of listings with at most one of each listing. Listing order is
// preserved by first occurrence in initial slice.
func UniqueListings(rawListings []lstg.Listing) []lstg.Listing {
	// nothing to filter here...
	if len(rawListings) == 0 {
		return []lstg.Listing{}
	}

	keys := make(map[lstg.Listing]bool)
	cleanListings := []lstg.Listing{}

	for _, listing := range rawListings {
		if _, found := keys[listing]; !found {
			keys[listing] = true
			cleanListings = append(cleanListings, listing)
		}
	}

	return cleanListings
}

func init() {
	importCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print imported listings to stdout.")
	rootCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
