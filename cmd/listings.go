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
	"github.com/spf13/viper"
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
		if searchYear >= viper.GetInt("search.min_year") {
			fmt.Println("Searching listings by year: ", searchYear)
		} else {
			fmt.Println("Searching listings by year: ANY")
		}

		if searchIssue >= viper.GetInt("search.min_issue") {
			fmt.Println("Searching listings by issue: ", searchIssue)
		} else {
			fmt.Println("Searching listings by issue: ANY")
		}

		if searchMember >= viper.GetInt("search.min_member") {
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
	listingsCmd.Flags().IntVarP(&searchYear, "year", "y", viper.GetInt("search.default_year"), "Search listings by LEX Issue year. Use '0' to search 'ANY'")
	listingsCmd.Flags().IntVarP(&searchIssue, "issue", "i", viper.GetInt("search.default_issue"), "Search listings by LEX Issue Number. Use '0' to search 'ANY'")
	listingsCmd.Flags().IntVarP(&searchMember, "member", "m", viper.GetInt("search.default_member"), "Search listings by member number. Use '0' to search 'ANY'")
}
