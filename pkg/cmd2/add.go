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
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
)

// RunAddListing adds a single listing to the datastore.
func RunAddListing(ll []Listing) (string, error) {
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

	if err != nil {
		log.Error("Failed to save to db: ")
		return "", err
	}

	log.WithFields(log.Fields{
		"cmd": "listings.add",
	}).Info("Adding a listing.")

	for _, l := range ll {
		// copy loop variable so i can accurately reference it for saving
		listing := l
		err = dsManager.Store.Save(&listing)
		if err != nil {
			log.Error("Failed to save new listing.")
			return "", err
		}
	}

	return Render(ll), nil
}
