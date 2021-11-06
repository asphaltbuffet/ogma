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
	"io/ioutil"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestRunImportListings(t *testing.T) {
	type args struct {
		fp string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RunImportListings(tt.args.fp)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunImportListings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RunImportListings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImportListings(t *testing.T) { //nolint:funlen // do this later 2021-11-05 BL
	type args struct {
		fp string
	}
	tests := []struct {
		name    string
		args    args
		want    Listings
		wantErr bool
	}{
		{
			name: "no file",
			args: args{
				fp: "test/b.json",
			},
			want:    Listings{},
			wantErr: true,
		},
		{
			name: "single entry",
			args: args{
				fp: "test/s.json",
			},
			want: Listings{
				Listings: []Listing{
					{Volume: 2, IssueNumber: 55, Year: 2021, Season: "Spring", PageNumber: 1, IndexedCategory: "Art & Photography", IndexedMemberNumber: 2989, MemberExtension: "", IsInternational: false, IsReview: false, ListingText: "Fingerpainting exchange.", IsArt: false, IsFlagged: false},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appFS := afero.NewMemMapFs()

			// create test files and directories
			err := appFS.MkdirAll("test", 0o755)
			assert.NoError(t, err)

			err = afero.WriteFile(appFS, "test/s.json", []byte("{\n"+
				"\"listings\": [\n"+
				"{\n"+
				"\"volume\": 2,\n"+
				"\"issue\": 55,\n"+
				"\"year\": 2021,\n"+
				"\"season\": \"Spring\",\n"+
				"\"page\": 1,\n"+
				"\"category\": \"Art & Photography\",\n"+
				"\"member\": 2989,\n"+
				"\"alt\": \"\",\n"+
				"\"international\": false,\n"+
				"\"review\": false,\n"+
				"\"text\": \"Fingerpainting exchange.\",\n"+
				"\"art\": false,\n"+
				"\"flag\": false\n"+
				"}\n"+
				"]\n"+
				"}"), 0o644)
			assert.NoError(t, err)

			testFile, _ := appFS.Open(tt.args.fp)

			got, err := ImportListings(testFile)
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

func init() {
	log.SetOutput(ioutil.Discard)
}
