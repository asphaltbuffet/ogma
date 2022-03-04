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

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
)

// Members is a container for multiple mail objects.
type Members struct {
	Mails []Mail `json:"members"`
}

// Member contains relevant information for a member.
type Member struct {
	ID      int    `storm:"id,increment"`
	Number  int    `json:"reference"`
	Name    string `json:"sender"`
	Address string `json:"receiver"`
}

const memberCommandLongDesc = "The member command allows you to add a new member to the tracker with name and/or address information."

var memberColumnConfigs = []table.ColumnConfig{
	{
		Name:  "Number",
		Align: text.AlignCenter,
	},
	{
		Name:  "Name",
		Align: text.AlignLeft,
	},
	{
		Name:  "Address",
		Align: text.AlignCenter,
	},
}

func init() {
	rootCmd.AddCommand(NewMemberCmd())
}

// NewMemberCmd creates a member command.
func NewMemberCmd() *cobra.Command {
	// cmd represents the member command
	cmd := &cobra.Command{
		Use:     "member",
		Short:   "Tracks penpal information",
		Long:    memberCommandLongDesc,
		Example: `ogma member --number=1234 --name="John Smith" --address="123 Fake St, Fakeville, FS 12345"`,
		Run:     RunMemberCmd,
	}

	cmd.Flags().IntP("number", "i", 0, "Member number.")
	cmd.Flags().StringP("name", "n", "", "Member name.")
	cmd.Flags().StringP("address", "a", "", "Mailing address.")

	return cmd
}

// RunMemberCmd implements functionality of a mail command.
func RunMemberCmd(cmd *cobra.Command, args []string) {
	m, err := memberFromArgs(cmd)
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
			"args":    args,
		}).Error("failed to parse member info arguments: ", err)
		cmd.PrintErrln("invalid member info input: ", err)
		return
	}

	dsManager, err := datastore.New(viper.GetString("datastore.filename"))
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
		}).Error("failed to open datastore when adding new member info: ", err)
		cmd.PrintErrln("failed to open datastore when adding new member info: ", err)
		return
	}
	defer dsManager.Stop()

	err = dsManager.Save(&m)
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
		}).Error("unable to save member info: ", err)
		cmd.PrintErrln("Failed to save member info: ", err)
		return
	}

	log.WithFields(log.Fields{
		"number":  m.Number,
		"name":    m.Name,
		"address": m.Address,
	}).Info("added member info")

	cmd.Printf("Added member info.")
}

// memberFromArgs creates a new member object from command arguments.
func memberFromArgs(cmd *cobra.Command) (Member, error) {
	i, err := cmd.Flags().GetInt("member")
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
			"args":    cmd.Args,
		}).Error("failed to get member number argument")
		return Member{}, fmt.Errorf("failed to get member number argument: %w", err)
	}

	n, err := cmd.Flags().GetString("name")
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
			"args":    cmd.Args,
		}).Error("failed to get name argument")
		return Member{}, fmt.Errorf("failed to get member name argument: %w", err)
	}

	a, err := cmd.Flags().GetString("address")
	if err != nil {
		log.WithFields(log.Fields{
			"command": cmd.Name(),
			"args":    cmd.Args,
		}).Error("failed to get member address argument")
		return Member{}, fmt.Errorf("failed to get member address argument: %w", err)
	}

	m := Member{
		Number:  i,
		Name:    n,
		Address: a,
	}

	return m, nil
}

// RenderMember returns a pretty formatted member info as table.
func RenderMember(mm []Member, p bool) string {
	// empty string if there no information to render
	if len(mm) == 0 {
		return "No member information found."
	}

	mt := table.NewWriter()

	mt.SetTitle("Member Information:")

	mt.AppendHeader(table.Row{
		"Number",
		"Name",
		"Address",
	})

	for _, m := range mm {
		mt.AppendRow([]interface{}{
			m.Number,
			m.Name,
			m.Address,
		})
	}

	mt.SetColumnConfigs(memberColumnConfigs)

	if p {
		mt.SetStyle(table.StyleColoredBright)
	}
	log.WithFields(log.Fields{
		"is_pretty": p,
	}).Debug("set rendering style")

	return mt.Render()
}
