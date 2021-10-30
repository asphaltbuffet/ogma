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

// listingsCmd represents the listings command.
var listingsCmd = &cobra.Command{
	Use:   "listings",
	Short: "Access listings functionality.",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

var searchListingsCmd = &cobra.Command{
	Use:   "search",
	Short: "Returns all listing information based on search criteria.",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunSearchListings(cmd)
	},
}

// RunSearchListings performs action associated with listings application command.
func RunSearchListings(cmd *cobra.Command) error {
	year, err := cmd.Flags().GetInt("year")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "year",
		}).Error("Invalid flag.")
		return err
	}

	issue, err := cmd.Flags().GetInt("issue")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "issue",
		}).Error("Invalid flag.")
		return err
	}

	member, err := cmd.Flags().GetInt("member")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "member",
		}).Error("Invalid flag.")
		return err
	}

	// log.WithFields(log.Fields{
	// 	"cmd":    "listings.search",
	// 	"year":   year,
	// 	"issue":  issue,
	// 	"member": member,
	// }).Info("Searching listings.")

	// TODO: SearchListings should return an error too.
	resultsCount := SearchListings(year, issue, member)
	cmd.Println("Found", resultsCount, "listings.")

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
	if y == 2021 { //nolint:gomnd // placeholder before db query is implemented
		count += 3
	} else {
		count += 0
	}

	if i == 56 { //nolint:gomnd // placeholder before db query is implemented
		count += 17
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

func init() {
	searchListingsCmd.Flags().IntP("year", "y", -1, "Search listings by LEX Issue year.")
	searchListingsCmd.Flags().IntP("issue", "i", -1, "Search listings by LEX Issue Number.")
	searchListingsCmd.Flags().IntP("member", "m", -1, "Search listings by member number.")

	listingsCmd.AddCommand(searchListingsCmd)
	rootCmd.AddCommand(listingsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listingsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listingsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
