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

// Package cmd contains all CLI commands implementations.
package cmd

import (
	"github.com/jedib0t/go-pretty/v6/table"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

// Render returns a pretty formatted listing as table.
func (l *Listing) Render() string {
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
		l.Volume,
		l.IssueNumber,
		l.Year,
		l.Season,
		l.PageNumber,
		l.IndexedCategory,
		l.IndexedMemberNumber,
		l.IsInternational,
		l.IsReview,
		l.ListingText,
		l.IsArt,
		l.IsFlagged,
	})
	lt.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:     "Text",
			WidthMax: viper.GetInt("defaults.max_column"),
		},
	})

	return lt.Render()
}

// ParseInputForListing creates a new listing from command flags.
func ParseInputForListing(c *cobra.Command) (l Listing, err error) { //nolint:funlen // TODO: refactor later 2021-11-02 BL
	volume, err := c.Flags().GetInt("volume")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "volume",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	lex, err := c.Flags().GetInt("lex")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "lex",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	year, err := c.Flags().GetInt("year")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "year",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	season, err := c.Flags().GetString("season")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "season",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	page, err := c.Flags().GetInt("page")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "page",
		}).Error("Invalid flag.")
		return Listing{}, err
	}
	category, err := c.Flags().GetString("category")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "category",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	member, err := c.Flags().GetInt("member")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "member",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	international, err := c.Flags().GetBool("international")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "international",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	review, err := c.Flags().GetBool("review")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "review",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	text, err := c.Flags().GetString("text")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "year",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	art, err := c.Flags().GetBool("art")
	if err != nil {
		log.WithFields(log.Fields{
			"flag": "art",
		}).Error("Invalid flag.")
		return Listing{}, err
	}

	flag, err := c.Flags().GetBool("flag")
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
