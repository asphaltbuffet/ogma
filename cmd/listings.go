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
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
)

// A Listing contains relevant information for LEX listings.
type Listing struct {
	ID                  int `storm:"id,increment"`
	Volume              int
	IssueNumber         int
	Year                int
	Season              string
	PageNumber          int
	IndexedCategory     string `storm:"index"`
	IndexedMemberNumber int    `storm:"index"`
	MemberExtension     string
	IsInternational     bool
	IsReview            bool
	ListingText         string
	IsArt               bool
	IsFlagged           bool
}

// listingsCmd represents the listings command.
var listingsCmd = &cobra.Command{
	Use:   "listings",
	Short: "Access listings functionality.",
	Long:  ``,
}

var searchListingsCmd = &cobra.Command{
	Use:   "search",
	Short: "Returns all listing information based on search criteria.",
	Long:  ``,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunSearchListings(cmd)
	},
}

var addListingCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a single listing.",
	Long:  ``,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunAddListing(cmd)
	},
}

var importListingsCmd = &cobra.Command{
	Use:   "import",
	Short: "Import listings from a file.",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunImportListings(cmd)
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

	log.WithFields(log.Fields{
		"cmd":    "listings.search",
		"year":   year,
		"issue":  issue,
		"member": member,
	}).Info("Searching listings.")

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

// ParseInputForListing creates a new listing from command flags.
func ParseInputForListing(cmd *cobra.Command) (l Listing, err error) { //nolint:funlen // TODO: refactor later 2021-11-02 BL
	volume, err := cmd.Flags().GetInt("volume")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "volume",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	lex, err := cmd.Flags().GetInt("lex")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "lex",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	year, err := cmd.Flags().GetInt("year")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "year",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	season, err := cmd.Flags().GetString("season")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "season",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	page, err := cmd.Flags().GetInt("page")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "page",
		}).Error("Invalid flag.")
		return Listing{}, err
	}
	category, err := cmd.Flags().GetString("category")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "category",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	member, err := cmd.Flags().GetInt("member")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "member",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	international, err := cmd.Flags().GetBool("international")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "international",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	review, err := cmd.Flags().GetBool("review")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "review",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	text, err := cmd.Flags().GetString("text")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "year",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	art, err := cmd.Flags().GetBool("art")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "art",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	flag, err := cmd.Flags().GetBool("flag")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "flag",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	return Listing{
		Volume:              volume,
		IssueNumber:         lex,
		Year:                year,
		Season:              season,
		PageNumber:          page,
		IndexedCategory:     category,
		IndexedMemberNumber: member,
		MemberExtension:     "", // TODO: add support for member extensions
		IsInternational:     international,
		IsReview:            review,
		ListingText:         text,
		IsArt:               art,
		IsFlagged:           flag,
	}, nil
}

// RunAddListing adds a single listing to the datastore.
func RunAddListing(cmd *cobra.Command) error { //nolint:funlen // TODO: refactor later 2021-10-31 BL
	dsManager, err := datastore.New(viper.GetString("datastore.filename"))
	if err != nil {
		log.Fatal("Datastore manager failure.")
	}

	defer dsManager.Stop()
	if err != nil {
		log.Error("Failed to save to db: ")
	}

	newListing, err := ParseInputForListing(cmd)
	if err != nil {
		log.Error("Failed to save to db: ")
		return err
	}

	// log.WithFields(log.Fields{
	// 	"cmd":    "listings.add",
	// 	"year":   year,
	// 	"issue":  lex,
	// 	"member": member,
	// }).Info("Adding a listing.")

	err = dsManager.Store.Save(&newListing)
	if err != nil {
		log.Error("Failed to save new listing.")
	}
	// TODO: add formatted listing as output.
	cmd.Println("Added a listing.")

	lt := table.NewWriter()
	lt.AppendHeader(table.Row{
		"Volume",
		"Issue",
		"Year",
		"Season",
		"Page",
		"Category",
		"Member",
		"International",
		"Review",
		"Text",
		"Sketch",
		"Flagged",
	})
	lt.AppendRow([]interface{}{
		newListing.Volume,
		newListing.IssueNumber,
		newListing.Year,
		newListing.Season,
		newListing.PageNumber,
		newListing.IndexedCategory,
		newListing.IndexedMemberNumber,
		newListing.IsInternational,
		newListing.IsReview,
		newListing.ListingText,
		newListing.IsArt,
		newListing.IsFlagged,
	})
	lt.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:     "Text",
			WidthMax: viper.GetInt("defaults.max_column"),
		},
	})

	// TODO: wrap output styles in a new flag
	// lt.SetStyle(table.StyleColoredBright)
	cmd.Println(lt.Render())

	// Actually return error if add was successful.
	return nil
}

// RunImportListings adds one to many listings to the datastore from a file.
func RunImportListings(cmd *cobra.Command) error {
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

func init() {
	importListingsCmd.Flags().Bool("verbose", false, "Print imported listings to stdout.")

	addListingCmd.Flags().IntP("volume", "v", -1, "Volume containing listing entry.")
	addListingCmd.Flags().IntP("lex", "l", viper.GetInt("defaults.issue"), "LEX issue containing listing entry.")
	addListingCmd.Flags().IntP("year", "y", time.Now().Year(), "Year of listing entry..")
	addListingCmd.Flags().StringP("season", "s", "", "Season of listing entry.")
	addListingCmd.Flags().IntP("page", "p", -1, "Page number of listing entry.")
	addListingCmd.Flags().StringP("category", "c", "", "Category of listing entry.")
	addListingCmd.Flags().IntP("member", "m", -1, "Member number of listing entry.")
	addListingCmd.Flags().BoolP("international", "i", false, "Is international postage required?")
	addListingCmd.Flags().BoolP("review", "r", false, "Is this a book review listing entry?")
	addListingCmd.Flags().StringP("text", "t", "", "Text of listing entry.")
	addListingCmd.Flags().BoolP("art", "a", false, "Is this a sketch listing entry?")
	addListingCmd.Flags().BoolP("flag", "f", false, "Has this listing entry been flagged?")
	listingsCmd.AddCommand(addListingCmd)

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
