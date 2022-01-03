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

// Package lstg contains all CLI commands implementations.
package lstg

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/jonreiter/govader"
	log "github.com/sirupsen/logrus"
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

var analyzer = govader.NewSentimentIntensityAnalyzer()

var listingColumnConfigs = []table.ColumnConfig{
	{
		Name:             "Text",
		WidthMax:         80, //nolint:gomnd // using viper fails unit tests. TODO: Fix this 2021-11-04 BL
		WidthMaxEnforcer: text.WrapSoft,
	},
	{
		Name:  "Member",
		Align: text.AlignRight,
		// AutoMerge: true, // doesn't look right without row separators
	},
	{
		Name:  "International",
		Align: text.AlignCenter,
	},
	{
		Name:  "Review",
		Align: text.AlignCenter,
	},
	{
		Name:  "Sketch",
		Align: text.AlignCenter,
	},
	{
		Name:  "Flagged",
		Align: text.AlignCenter,
	},
}

// RenderListings returns a pretty formatted listing as table.
func RenderListings(ll []Listing, p bool) string {
	// empty string if there are no listings to render
	if len(ll) == 0 {
		return "No LEX listings found."
	}

	lt := table.NewWriter()

	lt.SetTitle("LEX Issue Matches:")

	lt.AppendHeader(table.Row{
		"ID",
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
		"Sentiment",
	})

	for _, l := range ll {
		lt.AppendRow([]interface{}{
			l.ID,
			l.Volume,
			l.IssueNumber,
			l.Year,
			l.Season,
			l.PageNumber,
			l.IndexedCategory,
			fmt.Sprint(l.IndexedMemberNumber) + l.MemberExtension,
			convertBool(l.IsInternational),
			convertBool(l.IsReview),
			l.ListingText,
			convertBool(l.IsArt),
			convertBool(l.IsFlagged),
			fmt.Sprintf("%.2f", l.calcSentiment()),
		})
		// lt.AppendSeparator() // disabling for now, it makes it look messy - may want it configurable at run-time
	}

	lt.SetColumnConfigs(listingColumnConfigs)

	if p {
		lt.SetStyle(table.StyleColoredBright)
	}
	log.WithFields(log.Fields{
		"is_pretty": p,
	}).Debug("set rendering style")

	return lt.Render()
}

// convertBool strips out 'false' values for easier reading.
func convertBool(b bool) string {
	if b {
		return "✔"
	}
	return ""
}

func (l *Listing) calcSentiment() float64 {
	sentiment := analyzer.PolarityScores(l.ListingText)

	log.WithFields(log.Fields{
		"positive": sentiment.Positive,
		"negative": sentiment.Negative,
		"neutral":  sentiment.Neutral,
		"compound": sentiment.Compound,
	}).Debugf("sentiment analysis of listing id: %d", l.ID)

	return sentiment.Compound
}
