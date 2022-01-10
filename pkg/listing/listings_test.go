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

package lstg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	lstg "github.com/asphaltbuffet/ogma/pkg/listing"
)

func TestRenderListings(t *testing.T) {
	type args struct {
		ll []lstg.Listing
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty listing",
			args: args{
				[]lstg.Listing{},
			},
			want: `No LEX listings found.`,
		},
		{
			name: "single listing",
			args: args{
				[]lstg.Listing{
					{Volume: 2, IssueNumber: 55, Year: 2021, Season: "Spring", PageNumber: 1, IndexedCategory: "Art & Photography", IndexedMemberNumber: 2989, MemberExtension: "", IsInternational: false, IsReview: false, ListingText: "Fingerpainting exchange.", IsArt: false, IsFlagged: false},
				},
			},
			want: "+------------------------------------------------------------------------------------------------------------------------------------------------------------+\n| LEX Issue Matches:                                                                                                                                         |\n+----+--------+-------+------+--------+------+-------------------+--------+---------------+--------+--------------------------+--------+---------+-----------+\n| ID | VOLUME | ISSUE | YEAR | SEASON | PAGE | CATEGORY          | MEMBER | INTERNATIONAL | REVIEW | TEXT                     | SKETCH | FLAGGED | SENTIMENT |\n+----+--------+-------+------+--------+------+-------------------+--------+---------------+--------+--------------------------+--------+---------+-----------+\n|  0 |      2 |    55 | 2021 | Spring |    1 | Art & Photography |   2989 |               |        | Fingerpainting exchange. |        |         | 0.00      |\n+----+--------+-------+------+--------+------+-------------------+--------+---------------+--------+--------------------------+--------+---------+-----------+",
		},
		{
			name: "multiple listings",
			args: args{
				[]lstg.Listing{
					{Volume: 2, IssueNumber: 55, Year: 2021, Season: "Spring", PageNumber: 1, IndexedCategory: "Art & Photography", IndexedMemberNumber: 2989, MemberExtension: "", IsInternational: false, IsReview: false, ListingText: "Fingerpainting exchange.", IsArt: false, IsFlagged: false},
					{Volume: 2, IssueNumber: 56, Year: 2021, Season: "Spring", PageNumber: 2, IndexedCategory: "Crafts", IndexedMemberNumber: 12784, MemberExtension: "", IsInternational: true, IsReview: false, ListingText: "", IsArt: true, IsFlagged: false},
					{Volume: 2, IssueNumber: 56, Year: 2021, Season: "Spring", PageNumber: 2, IndexedCategory: "Creative Writing", IndexedMemberNumber: 11062, MemberExtension: "", IsInternational: false, IsReview: false, ListingText: "Writer's workshop zine.", IsArt: false, IsFlagged: true},
					{Volume: 2, IssueNumber: 56, Year: 2021, Season: "Spring", PageNumber: 2, IndexedCategory: "Creative Writing", IndexedMemberNumber: 7214, MemberExtension: "", IsInternational: false, IsReview: true, ListingText: "_Crimson Letters: Voices from Death Row_ consists of 30 compelling essays written in the prisoners' own words, offering stories of brutal beatings inside Juvenile Hall, botched suicide attempts, the terror of the first night on Death Row, the pain of goodbye as a friend is led to  execution, and the small acts of humanity that keep hope alive for men living in the shadow of death.", IsArt: false, IsFlagged: true},
				},
			},
			want: "+--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------+\n| LEX Issue Matches:                                                                                                                                                                                                 |\n+----+--------+-------+------+--------+------+-------------------+--------+---------------+--------+----------------------------------------------------------------------------------+--------+---------+-----------+\n| ID | VOLUME | ISSUE | YEAR | SEASON | PAGE | CATEGORY          | MEMBER | INTERNATIONAL | REVIEW | TEXT                                                                             | SKETCH | FLAGGED | SENTIMENT |\n+----+--------+-------+------+--------+------+-------------------+--------+---------------+--------+----------------------------------------------------------------------------------+--------+---------+-----------+\n|  0 |      2 |    55 | 2021 | Spring |    1 | Art & Photography |   2989 |               |        | Fingerpainting exchange.                                                         |        |         | 0.00      |\n|  0 |      2 |    56 | 2021 | Spring |    2 | Crafts            |  12784 |       ✔       |        |                                                                                  |    ✔   |         | 0.00      |\n|  0 |      2 |    56 | 2021 | Spring |    2 | Creative Writing  |  11062 |               |        | Writer's workshop zine.                                                          |        |    ✔    | 0.00      |\n|  0 |      2 |    56 | 2021 | Spring |    2 | Creative Writing  |   7214 |               |    ✔   | _Crimson Letters: Voices from Death Row_ consists of 30 compelling essays        |        |    ✔    | -0.97     |\n|    |        |       |      |        |      |                   |        |               |        | written in the prisoners' own words, offering stories of brutal beatings inside  |        |         |           |\n|    |        |       |      |        |      |                   |        |               |        | Juvenile Hall, botched suicide attempts, the terror of the first night on Death  |        |         |           |\n|    |        |       |      |        |      |                   |        |               |        | Row, the pain of goodbye as a friend is led to execution, and the small acts of  |        |         |           |\n|    |        |       |      |        |      |                   |        |               |        | humanity that keep hope alive for men living in the shadow of death.             |        |         |           |\n+----+--------+-------+------+--------+------+-------------------+--------+---------------+--------+----------------------------------------------------------------------------------+--------+---------+-----------+",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lstg.RenderListings(tt.args.ll, false)
			assert.Equal(t, tt.want, got)
		})
	}
}
