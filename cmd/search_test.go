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
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/asphaltbuffet/ogma/cmd"
	"github.com/asphaltbuffet/ogma/pkg/datastore"
	lstg "github.com/asphaltbuffet/ogma/pkg/listing"
)

func TestNewSearchCmd(t *testing.T) {
	got := cmd.NewSearchCmd()

	assert.Equal(t, "search", got.Name())
	// assert.Equal(t, "Returns all listing information based on search criteria.", got.Short)
	assert.True(t, got.Runnable())
}

func TestRunSearchCmd(t *testing.T) {
	m, dsFile := initDatastoreManager(t)
	m.Stop()

	defer func() {
		require.NoError(t, os.RemoveAll("test/"))
	}()

	tests := []struct {
		name      string
		args      []string
		datastore string
		assertion assert.ErrorAssertionFunc
		want      string
	}{
		{
			name:      "invalid datastore",
			args:      []string{"1234"},
			datastore: "notafile.db",
			assertion: assert.NoError,
			want:      "error opening datastore:  error accessing datastore file:",
		},
		{
			name:      "invalid search - too many parameters",
			args:      []string{"1234", "5678"},
			datastore: dsFile,
			assertion: assert.Error,
			want:      "Error: requires a single member number",
		},
		{
			name:      "invalid search - alphanumeric",
			args:      []string{"1234a"},
			datastore: dsFile,
			assertion: assert.Error,
			want:      "Error: invalid member number: strconv.Atoi: parsing \"1234a\": invalid syntax",
		},
		{
			name:      "with listings, no correspondence",
			args:      []string{"1234"},
			datastore: dsFile,
			assertion: assert.NoError,
			want:      "\n+---------------------------------------------------------------------------------------------------------------------------------------------------------+\n| LEX Issue Matches:                                                                                                                                      |\n+----+--------+-------+------+---------+------+----------+--------+---------------+--------+-------------------------------------------+--------+---------+\n| ID | VOLUME | ISSUE | YEAR | SEASON  | PAGE | CATEGORY | MEMBER | INTERNATIONAL | REVIEW | TEXT                                      | SKETCH | FLAGGED |\n+----+--------+-------+------+---------+------+----------+--------+---------------+--------+-------------------------------------------+--------+---------+\n|  1 |      1 |     1 | 1986 | Mollit  |    1 | Pariatur |   1234 |               |        | Esse Lorem do nulla sunt mollit nulla in. |        |    ✔    |\n|  2 |      1 |     1 | 1986 | Eiusmod |    2 | Commodo  |  1234B |               |        | Magna officia anim dolore enim.           |        |    ✔    |\n+----+--------+-------+------+---------+------+----------+--------+---------------+--------+-------------------------------------------+--------+---------+\n\n+------------------------------------------------------+\n| Correspondence Matches:                              |\n+-----------+--------+----------+------------+---------+\n| REFERENCE | SENDER | RECEIVER | DATE       | LINK    |\n+-----------+--------+----------+------------+---------+\n| 123d5f    |     55 |     1234 | 1986-04-01 |    L1   |\n| b12cd3    |   1234 |       55 | 1986-05-16 | M123d5f |\n| 6beef9    |   1234 |      666 | 2021-03-15 |         |\n+-----------+--------+----------+------------+---------+\n",
		},
		{
			name:      "no listings, with correspondence",
			args:      []string{"1234"},
			datastore: dsFile,
			assertion: assert.NoError,
			want:      "\n+---------------------------------------------------------------------------------------------------------------------------------------------------------+\n| LEX Issue Matches:                                                                                                                                      |\n+----+--------+-------+------+---------+------+----------+--------+---------------+--------+-------------------------------------------+--------+---------+\n| ID | VOLUME | ISSUE | YEAR | SEASON  | PAGE | CATEGORY | MEMBER | INTERNATIONAL | REVIEW | TEXT                                      | SKETCH | FLAGGED |\n+----+--------+-------+------+---------+------+----------+--------+---------------+--------+-------------------------------------------+--------+---------+\n|  1 |      1 |     1 | 1986 | Mollit  |    1 | Pariatur |   1234 |               |        | Esse Lorem do nulla sunt mollit nulla in. |        |    ✔    |\n|  2 |      1 |     1 | 1986 | Eiusmod |    2 | Commodo  |  1234B |               |        | Magna officia anim dolore enim.           |        |    ✔    |\n+----+--------+-------+------+---------+------+----------+--------+---------------+--------+-------------------------------------------+--------+---------+\n\n+------------------------------------------------------+\n| Correspondence Matches:                              |\n+-----------+--------+----------+------------+---------+\n| REFERENCE | SENDER | RECEIVER | DATE       | LINK    |\n+-----------+--------+----------+------------+---------+\n| 123d5f    |     55 |     1234 | 1986-04-01 |    L1   |\n| b12cd3    |   1234 |       55 | 1986-05-16 | M123d5f |\n| 6beef9    |   1234 |      666 | 2021-03-15 |         |\n+-----------+--------+----------+------------+---------+\n",
		},
		{
			name:      "no listings, no correspondence",
			args:      []string{"42"},
			datastore: dsFile,
			assertion: assert.NoError,
			want:      "\nNo LEX listings found.\n\nNo correspondences found.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("datastore.filename", tt.datastore)

			cmd := cmd.NewSearchCmd()
			b := bytes.NewBufferString("")
			cmd.SetOut(b)
			cmd.SetErr(b)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			tt.assertion(t, err)

			assert.Equal(t, tt.want, b.String()[:len(tt.want)], "unexpected output")
		})
	}
}

func initDatastoreManager(t *testing.T) (*datastore.Manager, string) {
	t.Helper()

	currentTime := time.Now()
	filename := fmt.Sprintf("test_%d.db", currentTime.Unix())
	manager, err := datastore.New(filename)
	require.NoError(t, err)

	appFS := afero.NewMemMapFs()

	// create test files and directories
	err = appFS.MkdirAll("test", 0o755)
	require.NoError(t, err)

	ltest := []lstg.Listing{
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
		{
			Volume:              1,
			IssueNumber:         1,
			Year:                1986,
			Season:              "Id",
			PageNumber:          3,
			IndexedCategory:     "Pariatur",
			IndexedMemberNumber: 5678,
			MemberExtension:     "",
			IsInternational:     false,
			IsReview:            false,
			ListingText:         "Velit cillum cillum ea officia nulla enim.",
			IsArt:               false,
			IsFlagged:           true,
		},
	}

	for _, record := range ltest {
		r := record
		_ = manager.Save(&r)
	}

	mtest := []cmd.Mail{
		{
			Ref:      "123d5f",
			Sender:   55,
			Receiver: 1234,
			Date:     "1986-04-01",
			Link:     "L1",
		},
		{
			Ref:      "b12cd3",
			Sender:   1234,
			Receiver: 55,
			Date:     "1986-05-16",
			Link:     "M123d5f",
		},
		{
			Ref:      "6beef9",
			Sender:   1234,
			Receiver: 666,
			Date:     "2021-03-15",
			Link:     "",
		},
	}

	for _, record := range mtest {
		r := record
		_ = manager.Save(&r)
	}

	return manager, filename
}

func init() {
	viper.GetViper().Set("search.max_results", 10)
}
