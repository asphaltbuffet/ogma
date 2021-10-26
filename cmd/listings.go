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
	"fmt"

	"github.com/spf13/cobra"
)

const (
	// MinYear - miniumum publication year.
	MinYear = 1986
	// MinIssue - miniumum publication issue.
	MinIssue = 1
	// MinMember - miniumum member number.
	MinMember = 1
)

var (
	searchYear   int
	searchIssue  int
	searchMember int
)

// listingsCmd represents the listings command.
var listingsCmd = &cobra.Command{
	Use:   "listings",
	Short: "Search listings by specified field.",
	Long: `Returns all listing information based on search criteria.
	
	Available search criteria:
	  Year
	  LEX Issue
	  Member Number
	  Category - NOT IMPLEMENTED
	  `,
	Run: func(cmd *cobra.Command, args []string) {
		if searchYear >= MinYear {
			fmt.Println("Searching listings by year: ", searchYear)
		} else {
			fmt.Println("Searching listings by year: ANY")
		}

		if searchIssue >= MinIssue {
			fmt.Println("Searching listings by issue: ", searchIssue)
		} else {
			fmt.Println("Searching listings by issue: ANY")
		}

		if searchMember >= MinMember {
			fmt.Println("Searching listings by member: ", searchMember)
		} else {
			fmt.Println("Searching listings by member: ANY")
		}
	},
}

func init() {
	rootCmd.AddCommand(listingsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listingsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listingsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listingsCmd.Flags().IntVarP(&searchYear, "year", "y", -1, "Search listings by LEX Issue year.")
	listingsCmd.Flags().IntVarP(&searchIssue, "issue", "i", -1, "Search listings by LEX Issue Number.")
	listingsCmd.Flags().IntVarP(&searchMember, "member", "m", -1, "Search listings by member number.")
}
