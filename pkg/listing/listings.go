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
	"errors"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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

var columnConfigs = []table.ColumnConfig{
	{
		Name:             "Text",
		WidthMax:         80, //nolint:gomnd // using viper fails unit tests. Fix this 2021-11-04 BL
		WidthMaxEnforcer: text.WrapSoft,
	},
	{
		Name:  "Member",
		Align: text.AlignRight,
		// AutoMerge: true, // doesn't look right without row separators
		// VAlign:    text.VAlignMiddle,
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

// GetStyle converts a []string into a valid rendering style.
func GetStyle(s []string) (table.Style, error) {
	// set default style as bright
	sa := "bright"

	// if too many args, just use the first one
	if len(s) > 0 {
		if len(s) > 1 {
			log.WithFields(log.Fields{
				"args": s,
			}).Warn("Too many styles passed in. Using first argument only.")
		}

		sa = s[0]
	}

	// check style
	switch sa {
	case "bright":
		return table.StyleColoredBright, nil
	case "light":
		return table.StyleLight, nil
	case "default":
		return table.StyleDefault, nil
	default:
		log.WithFields(log.Fields{
			"style": sa,
		}).Error("Invalid table style.")
		return table.StyleDefault, errors.New("invalid argument")
	}
}

// Render returns a pretty formatted listing as table.
func Render(ll []Listing, s ...string) string {
	// empty string if there are no listings to render
	if len(ll) == 0 {
		return ""
	}

	lt := table.NewWriter()

	// lt.SetTitle("Search results for Member #%d", ll[0].IndexedMemberNumber)

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
			ConvertBool(l.IsInternational),
			ConvertBool(l.IsReview),
			l.ListingText,
			ConvertBool(l.IsArt),
			ConvertBool(l.IsFlagged),
		})
		// lt.AppendSeparator() // disabling for now, it makes it look messy - may want it configurable at run-time
	}

	lt.SetColumnConfigs(columnConfigs)

	style, _ := GetStyle(s)
	lt.SetStyle(style)

	return lt.Render()
}

// ConvertBool strips out 'false' values for easier reading.
func ConvertBool(b bool) string {
	if b {
		return "✔"
	}
	return ""
}