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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	dsmocks "github.com/asphaltbuffet/ogma/pkg/datastore/mocks"
)

func TestAddListing(t *testing.T) {
	type args struct {
		ll []Listing
	}
	tests := []struct {
		name      string
		args      args
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			// TODO: This may make more sense to return an error. 2021-11-08
			name: "empty",
			args: args{
				ll: []Listing{},
			},
			want:      "",
			assertion: assert.NoError,
		},
		{
			name: "single entry",
			args: args{
				ll: []Listing{
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				},
			},
			want:      "+--------+-------+------+---------+------+----------+--------+---------------+--------+----------------------------------------------------------------------+--------+---------+\n| VOLUME | ISSUE | YEAR | SEASON  | PAGE | CATEGORY | MEMBER | INTERNATIONAL | REVIEW | TEXT                                                                 | SKETCH | FLAGGED |\n+--------+-------+------+---------+------+----------+--------+---------------+--------+----------------------------------------------------------------------+--------+---------+\n|      1 |     1 | 1999 | Qui sit |    1 | Pariatur |  99999 | false         | false  | Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit. | false  | false   |\n+--------+-------+------+---------+------+----------+--------+---------------+--------+----------------------------------------------------------------------+--------+---------+",
			assertion: assert.NoError,
		},
		{
			name: "multiple unique",
			args: args{
				ll: []Listing{
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 2, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				},
			},
			want:      "+--------+-------+------+---------+------+----------+--------+---------------+--------+----------------------------------------------------------------------+--------+---------+\n| VOLUME | ISSUE | YEAR | SEASON  | PAGE | CATEGORY | MEMBER | INTERNATIONAL | REVIEW | TEXT                                                                 | SKETCH | FLAGGED |\n+--------+-------+------+---------+------+----------+--------+---------------+--------+----------------------------------------------------------------------+--------+---------+\n|      1 |     1 | 1999 | Qui sit |    1 | Pariatur |  99999 | false         | false  | Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit. | false  | false   |\n|      2 |     1 | 1999 | Qui sit |    1 | Pariatur |  99999 | false         | false  | Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit. | false  | false   |\n+--------+-------+------+---------+------+----------+--------+---------------+--------+----------------------------------------------------------------------+--------+---------+",
			assertion: assert.NoError,
		},
		{
			name: "duplicates - single unique",
			args: args{
				ll: []Listing{
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				},
			},
			want:      "+--------+-------+------+---------+------+----------+--------+---------------+--------+----------------------------------------------------------------------+--------+---------+\n| VOLUME | ISSUE | YEAR | SEASON  | PAGE | CATEGORY | MEMBER | INTERNATIONAL | REVIEW | TEXT                                                                 | SKETCH | FLAGGED |\n+--------+-------+------+---------+------+----------+--------+---------------+--------+----------------------------------------------------------------------+--------+---------+\n|      1 |     1 | 1999 | Qui sit |    1 | Pariatur |  99999 | false         | false  | Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit. | false  | false   |\n+--------+-------+------+---------+------+----------+--------+---------------+--------+----------------------------------------------------------------------+--------+---------+",
			assertion: assert.NoError,
		},
		{
			name: "duplicates - multiple unique",
			args: args{
				ll: []Listing{
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 2, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				},
			},
			want:      "+--------+-------+------+---------+------+----------+--------+---------------+--------+----------------------------------------------------------------------+--------+---------+\n| VOLUME | ISSUE | YEAR | SEASON  | PAGE | CATEGORY | MEMBER | INTERNATIONAL | REVIEW | TEXT                                                                 | SKETCH | FLAGGED |\n+--------+-------+------+---------+------+----------+--------+---------------+--------+----------------------------------------------------------------------+--------+---------+\n|      1 |     1 | 1999 | Qui sit |    1 | Pariatur |  99999 | false         | false  | Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit. | false  | false   |\n|      2 |     1 | 1999 | Qui sit |    1 | Pariatur |  99999 | false         | false  | Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit. | false  | false   |\n+--------+-------+------+---------+------+----------+--------+---------------+--------+----------------------------------------------------------------------+--------+---------+",
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDatastore := &dsmocks.Writer{}
			mockDatastore.On("Save", mock.Anything).Return(nil)

			got, err := AddListing(tt.args.ll, mockDatastore)
			tt.assertion(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUniqueListings(t *testing.T) {
	type args struct {
		ll []Listing
	}
	tests := []struct {
		name string
		args args
		want []Listing
	}{
		{
			name: "empty",
			args: args{
				ll: []Listing{},
			},
			want: []Listing{},
		},
		{
			name: "no duplicates",
			args: args{
				ll: []Listing{
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 2, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				},
			},
			want: []Listing{
				{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				{Volume: 2, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
			},
		},
		{
			name: "only duplicates",
			args: args{
				ll: []Listing{
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				},
			},
			want: []Listing{
				{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
			},
		},
		{
			name: "duplicates with unique",
			args: args{
				ll: []Listing{
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 2, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
					{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				},
			},
			want: []Listing{
				{Volume: 1, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
				{Volume: 2, IssueNumber: 1, Year: 1999, Season: "Qui sit", PageNumber: 1, IndexedCategory: "Pariatur", IndexedMemberNumber: 99999, MemberExtension: "A", IsInternational: false, IsReview: false, ListingText: "Laborum aliquip eiusmod quis Lorem cupidatat nulla magna elit velit.", IsArt: false, IsFlagged: false},
			},
		},
		{
			name: "multiple duplicates with unique",
			args: args{
				ll: []Listing{
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
			want: []Listing{
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
			if got := UniqueListings(tt.args.ll); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UniqueListings() = %v, want %v", got, tt.want)
			}
		})
	}
}
