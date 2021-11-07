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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
)

// RunImportListings adds one to many listings to the datastore from a file.
func RunImportListings(f io.Reader, d datastore.Writer) (string, error) {
	// verify that the import reader is valid
	if f == nil {
		return "", errors.New("import : import reader cannot be nil")
	}

	// we initialize our listings array
	var listings Listings

	// convert import file into a listings struct
	listings, err := ImportListings(f)
	if err != nil {
		log.WithFields(log.Fields{"cmd": "listings.import"}).Error("Failed to import listings.")
		return "", err
	}

	// datastore needs to add one listing at a time, walk through imported listings and save one by one
	for _, l := range listings.Listings {
		listing := l

		err = d.Save(&listing)
		if err != nil {
			log.WithFields(log.Fields{"cmd": "listings.import", "count": len(listings.Listings)}).Error("Failed datastore save.")
			return "", err
		}
		log.WithFields(log.Fields{"cmd": "listings.import", "listing": listing}).Debug("Imported listing.")
	}
	log.WithFields(log.Fields{"cmd": "listings.import", "count": len(listings.Listings)}).Info("Imported listings.")

	return Render(listings.Listings), nil
}

// ImportListings unmarshalls a json file into a Listings struct.
func ImportListings(j io.Reader) (Listings, error) {
	// verify that the parameter is valid
	if j == nil {
		return Listings{}, errors.New("import : parameter cannot be nil")
	}

	// read our opened jsonFile as a byte array.
	byteValue, err := afero.ReadAll(j)
	if err != nil {
		log.WithFields(log.Fields{
			"cmd": "listings.import",
		}).Error("Failed to unmarshall import file.")
		return Listings{}, err
	}

	var ll Listings

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	err = json.Unmarshal(byteValue, &ll)
	if err != nil {
		log.WithFields(log.Fields{
			"cmd": "listings.import",
		}).Error("Failed to unmarshall import file.")
		return Listings{}, err
	}

	return ll, nil
}
