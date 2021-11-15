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
	"testing"

	"github.com/spf13/cobra"

	"github.com/asphaltbuffet/ogma/cmd"
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

// func TestSearch(t *testing.T) {
// 	var mockFinder *dsmocks.Finder

// 	// type finderMock struct {
// 	// 	finder *dsmocks.Finder
// 	// }
// 	tests := []struct {
// 		name   string
// 		member int
// 		// on        func(*finderMock)
// 		want      []lstg.Listing
// 		assertion assert.ErrorAssertionFunc
// 	}{
// 		{
// 			name:   "error",
// 			member: 12345,
// 			// on: func(f *finderMock) {
// 			// 	f.finder.On("Find", "IndexedMemberNumber", mock.AnythingOfType("int"), mock.Anything).Return(nil)
// 			// },
// 			want:      []lstg.Listing{},
// 			assertion: assert.Error,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// fMock := &finderMock{
// 			// 	&dsmocks.Finder{},
// 			// }
// 			// if tt.on != nil {
// 			// 	tt.on(fMock)
// 			// }

// 			// var ml []cmd2.Listing

// 			mockFinder.On("Find", "IndexedMemberNumber", tt.member).Return(func(s string, i int) interface{} {
// 				return nil
// 			}, func(s string, i int) error {
// 				return errors.New("testing error")
// 			})

// 			// mockFinder.On("Find", "IndexedMemberNumber", tt.member, ml).Return(nil).Run(func(args mock.Arguments) {
// 			// 	arg, ok := args.Get(2).(*[]cmd2.Listing)
// 			// 	assert.False(t, ok)
// 			// 	arg = append(arg, cmd2.Listing{
// 			// 		ID:                  123,
// 			// 		Volume:              45,
// 			// 		IssueNumber:         6,
// 			// 		Year:                7890,
// 			// 		Season:              "abcdef",
// 			// 		PageNumber:          1,
// 			// 		IndexedCategory:     "Ghi Jk",
// 			// 		IndexedMemberNumber: 23456,
// 			// 		MemberExtension:     "L",
// 			// 		IsInternational:     false,
// 			// 		IsReview:            false,
// 			// 		ListingText:         "Et reprehenderit duis consequat incididunt laborum commodo labore.",
// 			// 		IsArt:               false,
// 			// 		IsFlagged:           false,
// 			// 	})
// 			// })

// 			got, err := cmd.Search(tt.member, mockFinder)
// 			tt.assertion(t, err)
// 			if err != nil {
// 				return
// 			}
// 			assert.ObjectsAreEqual(tt.want, got)
// 		})
// 	}
// }
