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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
	lstg "github.com/asphaltbuffet/ogma/pkg/listing"
)

// NewSearchCmd creates a mail command.
func NewSearchCmd() *cobra.Command {
	// cmd represents the mail command
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Returns all listing information based on search criteria.",
		// Long:  `TODO: Add longer description about 'search'.`,
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
	// member number is already validated but check before using anyway.
	member, err := strconv.Atoi(args[0])
	if err != nil {
		log.WithField("arg", args[0]).Errorf("invalid member number: %v", err)

		cmd.PrintErrf("invalid member number: %v\n", err)
		return
	}

	log.WithField("member", member).Debug("searching by member number")

	dsManager, err := datastore.New(viper.GetString("datastore.filename"))
	if err != nil {
		log.Errorf("failed to access datastore: %v", err)

		cmd.PrintErrf("failed to access datastore: %v\n", err)
		return
	}
	defer dsManager.Stop()

	ll, err := SearchListings(member, dsManager)
	if err != nil {
		log.WithField("member", member).Errorf("failed to search listings: %v", err)

		cmd.PrintErrf("failed to search listings: %v\n", err)
		return
	}

	mm, err := SearchMail(member, dsManager)
	if err != nil {
		log.WithField("member", member).Errorf("failed to search mail: %v", err)

		cmd.PrintErrf("failed to search mail: %v", err)
		return
	}

	p, err := cmd.Flags().GetBool("pretty")
	if err != nil {
		log.Errorf("unable to read 'pretty' flag: %v", err)
		p = false
	}

	cmd.Printf("\n%s\n", lstg.RenderListings(ll, p))

	cmd.Printf("\n%s\n", RenderMail(mm, p))
}

// SearchListings returns all records with a matching member number (ignores member extensions).
func SearchListings(member int, ds storm.Finder) ([]lstg.Listing, error) {
	searchResults := []lstg.Listing{}
	err := ds.Find("IndexedMemberNumber", member, &searchResults)
	if err != nil {
		if !errors.Is(err, storm.ErrNotFound) {
			log.WithField("member", member).Errorf("failed to query database for listings: %v", err)
			return nil, fmt.Errorf("failure to query database for listings: %w", err)
		}

		log.WithField("member", member).Debug("no listings found")
	}

	return searchResults, nil
}

// SearchMail returns all mail records with a matching member number.
func SearchMail(member int, ds storm.Finder) ([]Mail, error) {
	var senderResults []Mail

	err := ds.Find("Sender", member, &senderResults)
	if err != nil {
		if !errors.Is(err, storm.ErrNotFound) {
			log.WithField("member", member).Errorf("failed to query database for mail by sender: %v", err)
			return nil, fmt.Errorf("failure to query database for mail by sender: %w", err)
		}

		log.WithField("sender", member).Debug("no sender correspondence found")
	}

	var receiverResults []Mail
	err = ds.Find("Receiver", member, &receiverResults)
	if err != nil {
		if !errors.Is(err, storm.ErrNotFound) {
			log.WithField("member", member).Errorf("failed to query database for mail by receiver: %v", err)
			return nil, fmt.Errorf("failure to query database for mail by receiver: %w", err)
		}

		log.WithField("receiver", member).Debug("no receiver correspondence found")
	}

	return append(senderResults, receiverResults...), nil
}

func init() {
	rootCmd.AddCommand(NewSearchCmd())
}
