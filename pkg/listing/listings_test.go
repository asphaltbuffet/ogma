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

func TestRender(t *testing.T) {
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
			want: ``,
		},
		{
			name: "single listing",
			args: args{
				[]lstg.Listing{
					{Volume: 2, IssueNumber: 55, Year: 2021, Season: "Spring", PageNumber: 1, IndexedCategory: "Art & Photography", IndexedMemberNumber: 2989, MemberExtension: "", IsInternational: false, IsReview: false, ListingText: "Fingerpainting exchange.", IsArt: false, IsFlagged: false},
				},
			},
			want: "\x1b[106;30m ID \x1b[0m\x1b[106;30m VOLUME \x1b[0m\x1b[106;30m ISSUE \x1b[0m\x1b[106;30m YEAR \x1b[0m\x1b[106;30m SEASON \x1b[0m\x1b[106;30m PAGE \x1b[0m\x1b[106;30m CATEGORY          \x1b[0m\x1b[106;30m MEMBER \x1b[0m\x1b[106;30m INTERNATIONAL \x1b[0m\x1b[106;30m REVIEW \x1b[0m\x1b[106;30m TEXT                     \x1b[0m\x1b[106;30m SKETCH \x1b[0m\x1b[106;30m FLAGGED \x1b[0m\n\x1b[107;30m  0 \x1b[0m\x1b[107;30m      2 \x1b[0m\x1b[107;30m    55 \x1b[0m\x1b[107;30m 2021 \x1b[0m\x1b[107;30m Spring \x1b[0m\x1b[107;30m    1 \x1b[0m\x1b[107;30m Art & Photography \x1b[0m\x1b[107;30m   2989 \x1b[0m\x1b[107;30m               \x1b[0m\x1b[107;30m        \x1b[0m\x1b[107;30m Fingerpainting exchange. \x1b[0m\x1b[107;30m        \x1b[0m\x1b[107;30m         \x1b[0m",
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
			want: "\x1b[106;30m ID \x1b[0m\x1b[106;30m VOLUME \x1b[0m\x1b[106;30m ISSUE \x1b[0m\x1b[106;30m YEAR \x1b[0m\x1b[106;30m SEASON \x1b[0m\x1b[106;30m PAGE \x1b[0m\x1b[106;30m CATEGORY          \x1b[0m\x1b[106;30m MEMBER \x1b[0m\x1b[106;30m INTERNATIONAL \x1b[0m\x1b[106;30m REVIEW \x1b[0m\x1b[106;30m TEXT                                                                             \x1b[0m\x1b[106;30m SKETCH \x1b[0m\x1b[106;30m FLAGGED \x1b[0m\n\x1b[107;30m  0 \x1b[0m\x1b[107;30m      2 \x1b[0m\x1b[107;30m    55 \x1b[0m\x1b[107;30m 2021 \x1b[0m\x1b[107;30m Spring \x1b[0m\x1b[107;30m    1 \x1b[0m\x1b[107;30m Art & Photography \x1b[0m\x1b[107;30m   2989 \x1b[0m\x1b[107;30m               \x1b[0m\x1b[107;30m        \x1b[0m\x1b[107;30m Fingerpainting exchange.                                                         \x1b[0m\x1b[107;30m        \x1b[0m\x1b[107;30m         \x1b[0m\n\x1b[47;30m  0 \x1b[0m\x1b[47;30m      2 \x1b[0m\x1b[47;30m    56 \x1b[0m\x1b[47;30m 2021 \x1b[0m\x1b[47;30m Spring \x1b[0m\x1b[47;30m    2 \x1b[0m\x1b[47;30m Crafts            \x1b[0m\x1b[47;30m  12784 \x1b[0m\x1b[47;30m       ✔       \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m                                                                                  \x1b[0m\x1b[47;30m    ✔   \x1b[0m\x1b[47;30m         \x1b[0m\n\x1b[107;30m  0 \x1b[0m\x1b[107;30m      2 \x1b[0m\x1b[107;30m    56 \x1b[0m\x1b[107;30m 2021 \x1b[0m\x1b[107;30m Spring \x1b[0m\x1b[107;30m    2 \x1b[0m\x1b[107;30m Creative Writing  \x1b[0m\x1b[107;30m  11062 \x1b[0m\x1b[107;30m               \x1b[0m\x1b[107;30m        \x1b[0m\x1b[107;30m Writer's workshop zine.                                                          \x1b[0m\x1b[107;30m        \x1b[0m\x1b[107;30m    ✔    \x1b[0m\n\x1b[47;30m  0 \x1b[0m\x1b[47;30m      2 \x1b[0m\x1b[47;30m    56 \x1b[0m\x1b[47;30m 2021 \x1b[0m\x1b[47;30m Spring \x1b[0m\x1b[47;30m    2 \x1b[0m\x1b[47;30m Creative Writing  \x1b[0m\x1b[47;30m   7214 \x1b[0m\x1b[47;30m               \x1b[0m\x1b[47;30m    ✔   \x1b[0m\x1b[47;30m _Crimson Letters: Voices from Death Row_ consists of 30 compelling essays        \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m    ✔    \x1b[0m\n\x1b[47;30m    \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m       \x1b[0m\x1b[47;30m      \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m      \x1b[0m\x1b[47;30m                   \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m               \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m written in the prisoners' own words, offering stories of brutal beatings inside  \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m         \x1b[0m\n\x1b[47;30m    \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m       \x1b[0m\x1b[47;30m      \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m      \x1b[0m\x1b[47;30m                   \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m               \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m Juvenile Hall, botched suicide attempts, the terror of the first night on Death  \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m         \x1b[0m\n\x1b[47;30m    \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m       \x1b[0m\x1b[47;30m      \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m      \x1b[0m\x1b[47;30m                   \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m               \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m Row, the pain of goodbye as a friend is led to execution, and the small acts of  \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m         \x1b[0m\n\x1b[47;30m    \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m       \x1b[0m\x1b[47;30m      \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m      \x1b[0m\x1b[47;30m                   \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m               \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m humanity that keep hope alive for men living in the shadow of death.             \x1b[0m\x1b[47;30m        \x1b[0m\x1b[47;30m         \x1b[0m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lstg.Render(tt.args.ll)
			assert.Equal(t, tt.want, got)
		})
	}
}
