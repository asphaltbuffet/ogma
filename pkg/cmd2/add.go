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

	"github.com/asphaltbuffet/ogma/pkg/datastore"
)

// AddListing adds a single listing to the datastore.
func AddListing(ll []Listing, ds datastore.Writer) (string, error) {
	// if there's nothing to add, return quickly
	if len(ll) == 0 {
		return "", nil
	}

	cl := UniqueListings(ll)

	log.WithFields(log.Fields{
		"cmd":        "listings.add",
		"count":      len(cl),
		"duplicates": len(ll) - len(cl),
	}).Info("Adding listing(s).")

	for _, l := range cl {
		// copy loop variable so i can accurately reference it for saving
		listing := l
		err := ds.Save(&listing)
		if err != nil {
			log.Error("Failed to save new listing.")
			return "", err
		}
	}

	return Render(cl), nil
}

// UniqueListings returns the passed in slice of listings with at most one of each listing. Listing order is
// preserved by first occurrence in initial slice.
func UniqueListings(ll []Listing) []Listing {
	// nothing to filter here...
	if len(ll) == 0 {
		return []Listing{}
	}

	keys := make(map[Listing]bool)
	cl := []Listing{}

	for _, l := range ll {
		if _, found := keys[l]; !found {
			keys[l] = true
			cl = append(cl, l)
		}
	}

	return cl
}
