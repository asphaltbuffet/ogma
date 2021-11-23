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
	"io/ioutil"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/asphaltbuffet/ogma/cmd"
	"github.com/asphaltbuffet/ogma/mocks"
	lstg "github.com/asphaltbuffet/ogma/pkg/listing"
)

func TestRunImport(t *testing.T) {
	tests := []struct {
		name      string
		filepath  string
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "no file",
			filepath:  "",
			want:      "",
			assertion: assert.Error,
		},
		{
			name:      "single entry",
			filepath:  "test/s.json",
			want:      "Imported 1 record.\n",
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appFS := afero.NewMemMapFs()

			// create test files and directories
			err := appFS.MkdirAll("test", 0o755)
			assert.NoError(t, err)
			if tt.filepath != "" {
				err = afero.WriteFile(appFS, tt.filepath, []byte(`{
					"listings": [
					{
					"volume": 2,
					"issue": 55,
					"year": 2021,
					"season": "Spring",
					"page": 1,
					"category": "Art & Photography",
					"member": 2989,
					"alt": "B",
					"international": false,
					"review": false,
					"text": "Fingerpainting exchange.",
					"art": false,
					"flag": false
					}
					]
					}`), 0o644)
				assert.NoError(t, err)
			}

			testFile, err := appFS.Open(tt.filepath)
			assert.NoError(t, err)

			mockDatastore := &mocks.Writer{}
			mockDatastore.On("Save", mock.Anything).Return(nil)

			got, err := cmd.Import(testFile, mockDatastore)
			tt.assertion(t, err)

			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestImportInput(t *testing.T) {
	type args struct {
		fp string
	}
	tests := []struct {
		name    string
		args    args
		want    []lstg.Listing
		wantErr bool
	}{
		{
			name: "no file",
			args: args{
				fp: "test/b.json",
			},
			want:    []lstg.Listing{},
			wantErr: true,
		},
		{
			name: "single entry",
			args: args{
				fp: "test/s.json",
			},
			want: []lstg.Listing{
				{Volume: 2, IssueNumber: 55, Year: 2021, Season: "Spring", PageNumber: 1, IndexedCategory: "Art & Photography", IndexedMemberNumber: 2989, MemberExtension: "", IsInternational: false, IsReview: false, ListingText: "Fingerpainting exchange.", IsArt: false, IsFlagged: false},
			},
			wantErr: false,
		},
		{
			name: "invalid json",
			args: args{
				fp: "test/invalid.json",
			},
			want:    []lstg.Listing{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appFS := afero.NewMemMapFs()

			// create test files and directories
			err := appFS.MkdirAll("test", 0o755)
			assert.NoError(t, err)

			err = afero.WriteFile(appFS, "test/s.json", []byte(`{
				"listings": [
					{
						"volume": 2,
						"issue": 55,
						"year": 2021,
						"season": "Spring",
						"page": 1,
						"category": "Art & Photography",
						"member": 2989,
						"alt": "",
						"international": false,
						"review": false,
						"text": "Fingerpainting exchange.",
						"art": false,
						"flag": false
					}
				]
				}`), 0o644)
			assert.NoError(t, err)

			err = afero.WriteFile(appFS, "test/invalid.json", []byte(`{
				"listings": [
					{
						"volume": 2,
						"issue": 55,
						"year": 2021,
						"season": "Spring",
						"page": 1,
						"category": "Art & Photography",
						"member": 2989,
						"alt": "",
				]
				}`), 0o644)
			assert.NoError(t, err)

			testFile, _ := appFS.Open(tt.args.fp)

			got, err := cmd.ParseImportInput(testFile)
			if err != nil {
				assert.Truef(t, tt.wantErr, "ImportListings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImportListings() = %v, want %v", got, tt.want)
			}
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

func ExecuteCommandC(t *testing.T, root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	t.Helper()

	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func init() {
	log.SetOutput(ioutil.Discard)

	// Search config in application directory with name ".ogma" (without extension).
	viper.AddConfigPath("../")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".ogma")

	viper.AutomaticEnv() // read in environment variables that match
}
