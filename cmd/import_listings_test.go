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
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/asphaltbuffet/ogma/cmd"
	lstg "github.com/asphaltbuffet/ogma/pkg/listing"
)

func TestNewImportListingCmd(t *testing.T) {
	got := cmd.NewImportListingCmd()

	assert.Equal(t, "listings", got.Name())
	assert.Equal(t, "Bulk import listing records.", got.Short)
	assert.True(t, got.Runnable())
}

func TestRunImportListingsCmd(t *testing.T) {
	m, dbFilePath, appFS := setup(t)
	m.Stop()

	defer func() {
		require.NoError(t, appFS.RemoveAll("test/"))
	}()

	tests := []struct {
		name      string
		args      []string
		datastore string
		assertion assert.ErrorAssertionFunc
		want      string
	}{
		{
			name:      "single entry",
			args:      []string{"test/listing.json"},
			datastore: dbFilePath,
			want:      "Imported 1/1 listing records.\n",
			assertion: assert.NoError,
		},
		{
			name:      "listing import",
			args:      []string{"test/listings.json"},
			datastore: dbFilePath,
			assertion: assert.NoError,
			want:      "Imported 3/3 listing records.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("datastore.filename", dbFilePath)

			cmd := cmd.NewImportListingCmd()
			b := bytes.NewBufferString("")
			cmd.SetOut(b)
			cmd.SetArgs(tt.args)
			tt.assertion(t, cmd.Execute())
			out, err := io.ReadAll(b)
			require.NoError(t, err)
			assert.Equal(t, tt.want, string(out))
		})
	}
}

func TestUniqueListings(t *testing.T) {
	type args struct {
		ll []lstg.Listing
	}
	tests := []struct {
		name string
		args args
		want []lstg.Listing
	}{
		{
			name: "empty",
			args: args{
				ll: []lstg.Listing{},
			},
			want: []lstg.Listing{},
		},
		{
			name: "no duplicates",
			args: args{
				ll: []lstg.Listing{
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 2, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				},
			},
			want: []lstg.Listing{
				{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				{Volume: 2, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
			},
		},
		{
			name: "only duplicates",
			args: args{
				ll: []lstg.Listing{
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				},
			},
			want: []lstg.Listing{
				{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
			},
		},
		{
			name: "duplicates with unique",
			args: args{
				ll: []lstg.Listing{
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 2, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				},
			},
			want: []lstg.Listing{
				{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				{Volume: 2, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
			},
		},
		{
			name: "multiple duplicates with unique",
			args: args{
				ll: []lstg.Listing{
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 2, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 3, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 2, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 3, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 3, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 2, IssueNumber: 2, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 4, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				},
			},
			want: []lstg.Listing{
				{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				{Volume: 2, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				{Volume: 3, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				{Volume: 2, IssueNumber: 2, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				{Volume: 4, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmd.UniqueListings(tt.args.ll); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UniqueListings() = %v, want %v", got, tt.want)
			}
		})
	}
}
