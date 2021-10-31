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

// RunAddListing adds a single listing to the datastore.
func RunAddListing(cmd *cobra.Command) error { //nolint:funlen // TODO: refactor later 2021-10-31 BL
	dsManager, err := datastore.New("ogma.db")
	if err != nil {
		log.Fatal("Datastore manager failure.")
	}

	defer dsManager.Stop()
	if err != nil {
		log.Error("Failed to save to db: ")
	}

	volume, err := cmd.Flags().GetInt("volume")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "volume",
		}).Error("Invalid flag.")
		return err
	}

	lex, err := cmd.Flags().GetInt("lex")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "lex",
		}).Error("Invalid flag.")
		return err
	}

	year, err := cmd.Flags().GetInt("year")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "year",
		}).Error("Invalid flag.")
		return err
	}

	page, err := cmd.Flags().GetInt("page")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "page",
		}).Error("Invalid flag.")
		return err
	}
	category, err := cmd.Flags().GetString("category")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "category",
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

	international, err := cmd.Flags().GetBool("international")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "international",
		}).Error("Invalid flag.")
		return err
	}

	review, err := cmd.Flags().GetBool("review")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "review",
		}).Error("Invalid flag.")
		return err
	}

	text, err := cmd.Flags().GetString("text")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "year",
		}).Error("Invalid flag.")
		return err
	}

	sketch, err := cmd.Flags().GetBool("sketch")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "sketch",
		}).Error("Invalid flag.")
		return err
	}

	flag, err := cmd.Flags().GetBool("flag")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "flag",
		}).Error("Invalid flag.")
		return err
	}

	log.WithFields(log.Fields{
		"cmd":    "listings.add",
		"year":   year,
		"issue":  lex,
		"member": member,
	}).Info("Adding a listing.")

	newListing := Listing{
		Volume:              volume,
		IssueNumber:         lex,
		Year:                year,
		PageNumber:          page,
		IndexedCategory:     category,
		IndexedMemberNumber: member,
		MemberExtension:     "", // TODO: add support for member extensions
		IsInternational:     international,
		IsReview:            review,
		ListingText:         text,
		IsArt:               sketch,
		IsFlagged:           flag,
	}

	err = dsManager.Store.Save(&newListing)
	if err != nil {
		log.Error("Failed to save new listing.")
	}
	// TODO: add formatted listing as output.
	cmd.Println("Added a listing.")

	// Actually return error if add was successful.
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
	addListingCmd.Flags().IntP("volume", "v", -1, "Volume containing listing entry.")
	addListingCmd.Flags().IntP("lex", "l", viper.GetInt("defaults.issue"), "LEX issue containing listing entry.")
	addListingCmd.Flags().IntP("year", "y", time.Now().Year(), "Year of listing entry..")
	// addListingCmd.MarkFlagRequired("year") //nolint:errcheck,gosec // handled by cobra
	addListingCmd.Flags().IntP("page", "p", -1, "Page number of listing entry.")
	addListingCmd.Flags().StringP("category", "c", "", "Category of listing entry.")
	addListingCmd.Flags().IntP("member", "m", -1, "Member number of listing entry.")
	addListingCmd.Flags().BoolP("international", "i", false, "Is international postage required?")
	addListingCmd.Flags().BoolP("review", "r", false, "Is this a book review listing entry?")
	addListingCmd.Flags().StringP("text", "t", "", "Text of listing entry.")
	addListingCmd.Flags().BoolP("sketch", "s", false, "Is this a sketch listing entry?")
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
