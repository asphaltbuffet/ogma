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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/pkg/cmd"
)

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
	RunE: func(c *cobra.Command, args []string) error {
		return cmd.RunSearchListings(c)
	},
}

var addListingCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a single listing.",
	Long:  ``,
	Args:  cobra.NoArgs,
	RunE: func(c *cobra.Command, args []string) error {
		return cmd.RunAddListing(c)
	},
}

var importListingsCmd = &cobra.Command{
	Use:   "import",
	Short: "Import listings from a file.",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	RunE: func(c *cobra.Command, args []string) error {
		return cmd.RunImportListings(c)
	},
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
