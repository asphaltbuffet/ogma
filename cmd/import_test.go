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
	dsmocks "github.com/asphaltbuffet/ogma/pkg/datastore/mocks"
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
			want:      "Imported 1 record.\n\x1b[106;30m ID \x1b[0m\x1b[106;30m VOLUME \x1b[0m\x1b[106;30m ISSUE \x1b[0m\x1b[106;30m YEAR \x1b[0m\x1b[106;30m SEASON \x1b[0m\x1b[106;30m PAGE \x1b[0m\x1b[106;30m CATEGORY          \x1b[0m\x1b[106;30m MEMBER \x1b[0m\x1b[106;30m INTERNATIONAL \x1b[0m\x1b[106;30m REVIEW \x1b[0m\x1b[106;30m TEXT                     \x1b[0m\x1b[106;30m SKETCH \x1b[0m\x1b[106;30m FLAGGED \x1b[0m\n\x1b[107;30m  0 \x1b[0m\x1b[107;30m      2 \x1b[0m\x1b[107;30m    55 \x1b[0m\x1b[107;30m 2021 \x1b[0m\x1b[107;30m Spring \x1b[0m\x1b[107;30m    1 \x1b[0m\x1b[107;30m Art & Photography \x1b[0m\x1b[107;30m  2989B \x1b[0m\x1b[107;30m               \x1b[0m\x1b[107;30m        \x1b[0m\x1b[107;30m Fingerpainting exchange. \x1b[0m\x1b[107;30m        \x1b[0m\x1b[107;30m         \x1b[0m",
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

			mockDatastore := &dsmocks.Writer{}
			mockDatastore.On("Save", mock.Anything).Return(nil)

			got, err := cmd.Import(testFile, mockDatastore)
			tt.assertion(t, err)

			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// func TestRunImportListingsCmd(t *testing.T) {
// 	var buf bytes.Buffer

// 	// If a config file is found, read it in.
// 	err := viper.ReadInConfig()
// 	assert.NoError(t, err)

// 	// Put the test_db in same place as config file for testing
// 	fp := filepath.Dir(viper.ConfigFileUsed())
// 	testDsfile := fp + "/test_db.db"

// 	// Change datastore for testing
// 	viper.Set("datastore.filename", testDsfile)
// 	defer func() {
// 		err = os.Remove(testDsfile)
// 		assert.NoError(t, err)
// 	}()
// 	tests := []struct {
// 		name    string
// 		args    []string
// 		want    string
// 		wantErr bool
// 	}{
// 		// // this testing doesn't allow verifying cobra behavior so far
// 		// {
// 		// 	name:    "no args",
// 		// 	args:    []string{"listings", "import"},
// 		// 	want:    "",
// 		// 	wantErr: true,
// 		// },
// 		{
// 			name:    "missing file",
// 			args:    []string{"listings", "import", fp + "/bad_file.json"},
// 			want:    "",
// 			wantErr: true,
// 		},
// 		{
// 			name:    "good import",
// 			args:    []string{"listings", "import", fp + "/importSingle_test.json"},
// 			want:    "",
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ogmaCmd := &cobra.Command{Use: "ogma"}
// 			listingsCmd := &cobra.Command{
// 				Use: "listings",
// 				RunE: func(c *cobra.Command, args []string) error {
// 					return nil
// 				},
// 			}

// 			importListingsCmd := &cobra.Command{
// 				Use: "import",
// 				RunE: func(c *cobra.Command, args []string) error {
// 					ogmaCmd.SetOut(&buf)
// 					err := RunImportListingsCmd(c, args)
// 					if (err != nil) != tt.wantErr {
// 						t.Errorf("RunImportListingsCmd() error = %v, wantErr %v", err, tt.wantErr)
// 					}
// 					return nil
// 				},
// 			}

// 			importListingsCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print imported listings to stdout.")
// 			listingsCmd.AddCommand(importListingsCmd)
// 			listingsCmd.AddCommand(importListingsCmd)
// 			ogmaCmd.AddCommand(listingsCmd)

// 			c, out, err := ExecuteCommandC(t, ogmaCmd, tt.args...)
// 			assert.Emptyf(t, out, "Unexpected output: %v", out)
// 			assert.NoError(t, err)
// 			assert.Equal(t, tt.want, buf.String())
// 			assert.Equalf(t, "import", c.Name(), `Invalid command returned from ExecuteC: expected "import", got: %q`, c.Name())
// 		})
// 	}
// }

func TestImportInput(t *testing.T) {
	type args struct {
		fp string
	}
	tests := []struct {
		name    string
		args    args
		want    lstg.Listings
		wantErr bool
	}{
		{
			name: "no file",
			args: args{
				fp: "test/b.json",
			},
			want:    lstg.Listings{},
			wantErr: true,
		},
		{
			name: "single entry",
			args: args{
				fp: "test/s.json",
			},
			want: lstg.Listings{
				Listings: []lstg.Listing{
					{Volume: 2, IssueNumber: 55, Year: 2021, Season: "Spring", PageNumber: 1, IndexedCategory: "Art & Photography", IndexedMemberNumber: 2989, MemberExtension: "", IsInternational: false, IsReview: false, ListingText: "Fingerpainting exchange.", IsArt: false, IsFlagged: false},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid json",
			args: args{
				fp: "test/invalid.json",
			},
			want:    lstg.Listings{},
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
