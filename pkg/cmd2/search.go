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
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RunSearchListings performs action associated with listings application command.
func RunSearchListings(c *cobra.Command) error {
	year, err := c.Flags().GetInt("year")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "year",
		}).Error("Invalid flag.")
		return err
	}

	issue, err := c.Flags().GetInt("issue")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "issue",
		}).Error("Invalid flag.")
		return err
	}

	member, err := c.Flags().GetInt("member")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "member",
		}).Error("Invalid flag.")
		return err
	}

	log.WithFields(log.Fields{
		"cmd":    "listings.search",
		"year":   year,
		"issue":  issue,
		"member": member,
	}).Info("Searching listings.")

	// TODO: SearchListings should return an error too.
	resultsCount := SearchListings(year, issue, member)
	c.Println("Found", resultsCount, "listings.")

	if resultsCount > viper.GetInt("search.max_results") {
		log.WithFields(log.Fields{
			"cmd":          "listings.search",
			"resultsCount": resultsCount,
			"maxResults":   viper.GetInt("search.max_results"),
		}).Warn("Query return too large.")
		return errors.New("too many results")
	}

	return nil
}

// SearchListings queries listing db for matching listings.
// TODO: Actually interact with db so it's functional.
func SearchListings(y int, i int, m int) (count int) {
	count = 0
	if y == 2021 { //nolint:gomnd // placeholder before db query is implemented
		count += 3
	} else {
		count += 0
	}

	if i == 56 { //nolint:gomnd // placeholder before db query is implemented
		count += 7
	} else {
		count += 0
	}

	if m == 1000 { //nolint:gomnd // placeholder before db query is implemented
		count += 2
	} else {
		count += 0
	}

	return count
}
