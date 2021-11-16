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
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
	lstg "github.com/asphaltbuffet/ogma/pkg/listing"
)

// addCmd represents the add command.
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a single listing.",
	Long:  ``,
	Args:  cobra.NoArgs,
	RunE:  RunAddListingCmd,
}

var (
	volume        int
	lex           int
	year          int
	season        string
	page          int
	category      string
	member        int
	international bool
	review        bool
	text          string
	art           bool
	flag          bool
)

// RunAddListingCmd performs action associated with listings-add application command.
func RunAddListingCmd(c *cobra.Command, args []string) error {
	dsManager, err := datastore.New(viper.GetString("datastore.filename"))
	if err != nil {
		log.Error("Datastore manager failure.")
		return err
	}
	defer dsManager.Stop()

	out, err := AddListing([]lstg.Listing{
		{
			Volume:              volume,
			IssueNumber:         lex,
			Year:                year,
			Season:              season,
			PageNumber:          page,
			IndexedCategory:     category,
			IndexedMemberNumber: member,
			MemberExtension:     "",
			IsInternational:     international,
			IsReview:            review,
			ListingText:         text,
			IsArt:               art,
			IsFlagged:           flag,
		},
	}, dsManager)
	if err == nil {
		c.Println(out)
	}

	return err
}

// AddListing adds a single listing to the datastore.
func AddListing(newListings []lstg.Listing, ds datastore.Writer) (string, error) {
	// if there's nothing to add, return quickly
	if len(newListings) == 0 {
		return "", nil
	}

	cl := UniqueListings(newListings)

	log.WithFields(log.Fields{
		"cmd":        "listings.add",
		"count":      len(cl),
		"duplicates": len(newListings) - len(cl),
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

	return lstg.Render(cl), nil
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
	addCmd.Flags().IntVarP(&volume, "volume", "v", -1, "Volume containing listing entry.")
	addCmd.Flags().IntVarP(&lex, "lex", "l", viper.GetInt("defaults.issue"), "LEX issue containing listing entry.")
	addCmd.Flags().IntVarP(&year, "year", "y", time.Now().Year(), "Year of listing entry..")
	addCmd.Flags().StringVarP(&season, "season", "s", "", "Season of listing entry.")
	addCmd.Flags().IntVarP(&page, "page", "p", -1, "Page number of listing entry.")
	addCmd.Flags().StringVarP(&category, "category", "c", "", "Category of listing entry.")
	addCmd.Flags().IntVarP(&member, "member", "m", -1, "Member number of listing entry.")
	addCmd.Flags().BoolVarP(&international, "international", "i", false, "Is international postage required?")
	addCmd.Flags().BoolVarP(&review, "review", "r", false, "Is this a book review listing entry?")
	addCmd.Flags().StringVarP(&text, "text", "t", "", "Text of listing entry.")
	addCmd.Flags().BoolVarP(&art, "art", "a", false, "Is this a sketch listing entry?")
	addCmd.Flags().BoolVarP(&flag, "flag", "f", false, "Has this listing entry been flagged?")

	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
