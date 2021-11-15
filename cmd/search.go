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
	"strconv"

	"github.com/asdine/storm/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
	lstg "github.com/asphaltbuffet/ogma/pkg/listing"
)

// searchCmd represents the search command.
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Returns all listing information based on search criteria.",
	Long:  `TODO: Add longer description about 'search'.`,
	Args:  cobra.ExactArgs(1),
	RunE:  RunSearchCmd,
}

// RunSearchCmd performs action associated with listings application command.
func RunSearchCmd(c *cobra.Command, args []string) error {
	member, err := strconv.Atoi(args[0])
	if err != nil {
		log.WithFields(log.Fields{
			"arg": args[0],
		}).Error("Invalid argument.")
		return errors.New("argument must be an integer")
	}

	log.WithFields(log.Fields{
		"cmd":    "search",
		"member": member,
	}).Info("Searching listings.")

	dsManager, err := datastore.New(viper.GetString("datastore.filename"))
	if err != nil {
		log.Error("Datastore manager failure.")
		return err
	}
	defer dsManager.Stop()

	ll, err := Search(member, dsManager)
	if err != nil {
		log.WithFields(log.Fields{
			"member": member,
		}).Error("Unable to get search results.")
		return err
	}

	// c.Printf("Found %d listings.\n", len(ll))
	c.Printf(lstg.Render(ll))

	return nil
}

// Search returns all records with a matching member number (ignores member extensions).
func Search(member int, ds storm.Finder) ([]lstg.Listing, error) {
	var searchResults []lstg.Listing
	err := ds.Find("IndexedMemberNumber", member, &searchResults)
	if err != nil {
		log.WithFields(log.Fields{
			"cmd":            "search",
			"member":         member,
			"internal_error": err,
		}).Error("Search failure.")
		return nil, errors.New("search: query failed")
	}

	if len(searchResults) > viper.GetInt("search.max_results") {
		log.WithFields(log.Fields{
			"cmd":          "search",
			"resultsCount": len(searchResults),
			"maxResults":   viper.GetInt("search.max_results"),
		}).Warn("Query return too large.")
		return nil, errors.New("too many results")
	}

	return searchResults, nil
}

func init() {
	searchCmd.Flags().IntP("member", "m", -1, "Search listings by member number.")
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
