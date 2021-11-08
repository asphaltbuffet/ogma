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
	"github.com/jedib0t/go-pretty/v6/text"
)

// A Listings hold unmarshalled listing data for import.
type Listings struct {
	Listings []Listing `json:"listings"`
}

// A Listing contains relevant information for LEX listings.
type Listing struct {
	ID                  int    `storm:"id,increment"`
	Volume              int    `json:"volume"`
	IssueNumber         int    `json:"issue"`
	Year                int    `json:"year"`
	Season              string `json:"season"`
	PageNumber          int    `json:"page"`
	IndexedCategory     string `storm:"index" json:"category"`
	IndexedMemberNumber int    `storm:"index" json:"member"`
	MemberExtension     string `json:"alt"`
	IsInternational     bool   `json:"international"`
	IsReview            bool   `json:"review"`
	ListingText         string `json:"text"`
	IsArt               bool   `json:"art"`
	IsFlagged           bool   `json:"flag"`
}

// Render returns a pretty formatted listing as table.
func Render(ll []Listing) string {
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

	for _, l := range ll {
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
	}

	lt.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:             "Text",
			WidthMax:         80, //nolint:gomnd // using viper fails unit tests. Fix this 2021-11-04 BL
			WidthMaxEnforcer: text.WrapSoft,
		},
	})

	return lt.Render()
}
