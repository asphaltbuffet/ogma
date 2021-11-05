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
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
)

// A Listings hold unmarshalled listing data for import.
type Listings struct {
	Listings []Listing `json:"listings"`
}

// RunImportListings adds one to many listings to the datastore from a file.
func RunImportListings(fp string) (string, error) {
	dsManager, err := datastore.New(viper.GetString("datastore.filename"))
	if err != nil {
		log.Error("Datastore manager failure.")
		return "", err
	}

	defer dsManager.Stop()
	if err != nil {
		log.Error("Failed to save to db: ")
		return "", err
	}

	// we initialize our listings array
	var listings Listings

	listings, err = ImportListings(fp)
	if err != nil {
		log.WithFields(log.Fields{
			"cmd":        "listings.import",
			"importFile": fp,
		}).Error("Failed to import listings.")
		return "", err
	}

	log.WithFields(log.Fields{
		"cmd":   "listings.import",
		"count": len(listings.Listings),
	}).Info("Imported listings.")

	for _, l := range listings.Listings {
		listing := l

		err = dsManager.Store.Save(&listing)
		if err != nil {
			log.Error("Failed to save new listings.")
			return "", err
		}
	}

	return Render(listings.Listings), nil
}

// ImportListings unmarshalls a json file into a Listings struct.
func ImportListings(fp string) (Listings, error) {
	jsonFile, err := os.Open(filepath.Clean(fp))
	if err != nil {
		log.WithFields(log.Fields{
			"cmd":        "listings.import",
			"importFile": fp,
		}).Error("Failed to open import file.")
	}

	log.WithFields(log.Fields{
		"cmd":        "listings.import",
		"importFile": fp,
	}).Info("Successfully opened import file.")

	defer func() {
		err = jsonFile.Close()
		if err != nil {
			log.WithFields(log.Fields{
				"cmd":        "listings.import",
				"importFile": fp,
			}).Error("Failed to close import file.")
		}
	}()

	// read our opened jsonFile as a byte array.
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.WithFields(log.Fields{
			"cmd":        "listings.import",
			"importFile": fp,
		}).Error("Failed to unmarshall import file.")
		return Listings{}, err
	}

	var ll Listings

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &ll)
	if err != nil {
		log.WithFields(log.Fields{
			"cmd":        "listings.import",
			"importFile": fp,
		}).Error("Failed to unmarshall import file.")
		return Listings{}, err
	}

	return ll, nil
}
