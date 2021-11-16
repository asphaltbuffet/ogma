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

package cmd_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/ogma/cmd"
	"github.com/asphaltbuffet/ogma/pkg/datastore"
	lstg "github.com/asphaltbuffet/ogma/pkg/listing"
)

func TestRunSearchCmd(t *testing.T) {
	type args struct {
		c    *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cmd.RunSearchCmd(tt.args.c, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("RunSearchCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	m, dbFilePath, err := initDatastoreManager()
	assert.NoError(t, err)
	defer m.Stop()

	defer func() {
		err = os.Remove(dbFilePath)
		assert.NoError(t, err)
	}()

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Mollit",
			PageNumber:          1,
			IndexedCategory:     "Pariatur",
			IndexedMemberNumber: 1234,
			MemberExtension:     "",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Esse Lorem do nulla sunt mollit nulla in.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	_, err = cmd.AddListing([]lstg.Listing{
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Eiusmod",
			PageNumber:          2,
			IndexedCategory:     "Commodo",
			IndexedMemberNumber: 1234,
			MemberExtension:     "B",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Magna officia anim dolore enim.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}, m)
	assert.NoError(t, err)

	for i := 1; i <= 10; i++ {
		_, err = cmd.AddListing([]lstg.Listing{
			{
				Volume:              1,
				IssueNumber:         1,
				Year:                1986,
				Season:              "Id",
				PageNumber:          3,
				IndexedCategory:     "Consequat",
				IndexedMemberNumber: 5678,
				MemberExtension:     "",
				IsInternational:     false,
				IsReview:            false,
				ListingText:         "Velit cillum cillum ea officia nulla enim.",
				IsArt:               false,
				IsFlagged:           true,
			},
		}, m)
		assert.NoError(t, err)
	}

	type args struct {
		member int
	}
	tests := []struct {
		name    string
		args    args
		want    []lstg.Listing
		wantErr bool
	}{
		{
			name: "find multiple",
			args: args{
				member: 1234,
			},
			want: []lstg.Listing{
				{
					ID:                  1,
					Volume:              1,
					IssueNumber:         1,
					Year:                1986,
					Season:              "Mollit",
					PageNumber:          1,
					IndexedCategory:     "Pariatur",
					IndexedMemberNumber: 1234,
					MemberExtension:     "",
					IsInternational:     false,
					IsReview:            false,
					ListingText:         "Esse Lorem do nulla sunt mollit nulla in.",
					IsArt:               false,
					IsFlagged:           true,
				},
				{
					ID:                  2,
					Volume:              1,
					IssueNumber:         1,
					Year:                1986,
					Season:              "Eiusmod",
					PageNumber:          2,
					IndexedCategory:     "Commodo",
					IndexedMemberNumber: 1234,
					MemberExtension:     "B",
					IsInternational:     false,
					IsReview:            false,
					ListingText:         "Magna officia anim dolore enim.",
					IsArt:               false,
					IsFlagged:           true,
				},
			},
			wantErr: false,
		},
		{
			name: "no results",
			args: args{
				member: 1,
			},
			want:    nil,
			wantErr: true,
		},
		// TODO: checked against max_results in config
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.Search(tt.args.member, m)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func initDatastoreManager() (*datastore.Manager, string, error) {
	currentTime := time.Now()
	filename := fmt.Sprintf("test_%d.db", currentTime.Unix())
	manager, err := datastore.New(filename)

	return manager, filename, err
}

func init() {
	viper.GetViper().Set("search.max_results", 10)
}
