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
	"fmt"
	"strconv"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
	lstg "github.com/asphaltbuffet/ogma/pkg/listing"
)

const searchCommandLongDesc = "The search command queries all LEX ads by the member who placed the ad. It" +
	"also searches any mail records by sender and receiver for matches with provided member number.\n\n" +
	"By default, it displays results in a colored table output. This can be changed with the '--pretty=false' flag."

func init() {
	rootCmd.AddCommand(NewSearchCmd())
}

// NewSearchCmd creates a mail command.
func NewSearchCmd() *cobra.Command {
	// cmd represents the mail command
	cmd := &cobra.Command{
		Use:   "search [member number]",
		Short: "Search records by member number.",
		Long:  searchCommandLongDesc,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("requires a single member number")
			}

			if _, err := strconv.Atoi(args[0]); err != nil {
				return fmt.Errorf("invalid member number: %w", err)
			}

			return nil
		},
		Run: RunSearchCmd,
	}

	cmd.Flags().BoolP("pretty", "p", false, "Show prettier results.")

	return cmd
}

// RunSearchCmd performs action associated with listings application command.
func RunSearchCmd(cmd *cobra.Command, args []string) {
	// member number is already validated by cobra
	member, _ := strconv.Atoi(args[0])

	log.WithField("member", member).Debug("searching by member number")

	dsManager, err := datastore.Open(viper.GetString("datastore.filename"))
	if err != nil {
		log.Error("error opening datastore: ", err)

		cmd.Println("error opening datastore: ", err)
		return
	}
	defer dsManager.Stop()

	ll, err := searchListings(member, dsManager)
	if err != nil {
		log.WithField("member", member).Error("failed to search listings: ", err)

		cmd.PrintErrln("failed to search listings: ", err)
		return
	}

	mm, err := SearchMail(member, dsManager)
	if err != nil {
		log.WithField("member", member).Error("failed to search mail: ", err)

		cmd.PrintErrln("failed to search mail: ", err)
		return
	}

	p, err := cmd.Flags().GetBool("pretty")
	if err != nil {
		log.Error("unable to read 'pretty' flag: ", err)
		p = false
	}

	cmd.Printf("\n%s\n", lstg.RenderListings(ll, p))

	cmd.Printf("\n%s\n", RenderMail(mm, p))
}

// searchListings returns all records with a matching member number (ignores member extensions).
func searchListings(member int, ds storm.Finder) ([]lstg.Listing, error) {
	searchResults := []lstg.Listing{}

	err := ds.Find("IndexedMemberNumber", member, &searchResults)
	if err != nil {
		switch err {
		case storm.ErrNotFound:
			log.WithField("member", member).Debug("no listings found")
		default:
			log.WithField("member", member).Error("failed to query database for listings: ", err)
			return nil, fmt.Errorf("failure to query database for listings: %w", err)
		}
	}

	return searchResults, nil
}

// SearchMail returns all mail records with a matching member number.
func SearchMail(member int, ds storm.Finder) ([]Mail, error) {
	var searchResults []Mail

	mailQuery := ds.Select(q.Or(q.Eq("Sender", member), q.Eq("Receiver", member)))

	err := mailQuery.Find(&searchResults)
	if err != nil {
		switch err {
		case storm.ErrNotFound:
			log.WithField("sender", member).Debug("no sender correspondence found")
		default:
			log.WithField("member", member).Error("failed to query database for mail by sender: ", err)
			return nil, fmt.Errorf("failure to query database for mail by sender=%d: %w", member, err)
		}
	}

	return searchResults, nil
}
