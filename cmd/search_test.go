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
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/ogma/cmd"
	"github.com/asphaltbuffet/ogma/pkg/datastore"
	lstg "github.com/asphaltbuffet/ogma/pkg/listing"
)

func TestNewSearchCmd(t *testing.T) {
	got := cmd.NewSearchCmd()

	assert.Equal(t, "search", got.Name())
	assert.Equal(t, "Returns all listing information based on search criteria.", got.Short)
	assert.True(t, got.Runnable())
}

func TestRunSearchCmd(t *testing.T) {
	m, dbFilePath, err := initDatastoreManager()
	assert.NoError(t, err)
	m.Stop()

	defer func() {
		assert.NoError(t, os.Remove(dbFilePath))
	}()

	viper.Set("datastore.filename", dbFilePath)
	tests := []struct {
		name      string
		args      []string
		assertion assert.ErrorAssertionFunc
		want      string
	}{
		{
			name:      "with listings, no correspondence",
			args:      []string{"1234"},
			assertion: assert.NoError,
			want:      "\n+---------------------------------------------------------------------------------------------------------------------------------------------------------+\n| LEX Issue Matches:                                                                                                                                      |\n+----+--------+-------+------+---------+------+----------+--------+---------------+--------+-------------------------------------------+--------+---------+\n| ID | VOLUME | ISSUE | YEAR | SEASON  | PAGE | CATEGORY | MEMBER | INTERNATIONAL | REVIEW | TEXT                                      | SKETCH | FLAGGED |\n+----+--------+-------+------+---------+------+----------+--------+---------------+--------+-------------------------------------------+--------+---------+\n|  1 |      1 |     1 | 1986 | Mollit  |    1 | Pariatur |   1234 |               |        | Esse Lorem do nulla sunt mollit nulla in. |        |    ✔    |\n|  2 |      1 |     1 | 1986 | Eiusmod |    2 | Commodo  |  1234B |               |        | Magna officia anim dolore enim.           |        |    ✔    |\n+----+--------+-------+------+---------+------+----------+--------+---------------+--------+-------------------------------------------+--------+---------+\n\n+------------------------------------------------------+\n| Correspondence Matches:                              |\n+-----------+--------+----------+------------+---------+\n| REFERENCE | SENDER | RECEIVER | DATE       | LINK    |\n+-----------+--------+----------+------------+---------+\n| 123d5f    |     55 |     1234 | 1986-04-01 |    L1   |\n| b12cd3    |   1234 |       55 | 1986-05-16 | M123d5f |\n| 6beef9    |   1234 |      666 | 2021-03-15 |         |\n+-----------+--------+----------+------------+---------+\n",
		},
		{
			name:      "no listings, with correspondence",
			args:      []string{"1234"},
			assertion: assert.NoError,
			want:      "\n+---------------------------------------------------------------------------------------------------------------------------------------------------------+\n| LEX Issue Matches:                                                                                                                                      |\n+----+--------+-------+------+---------+------+----------+--------+---------------+--------+-------------------------------------------+--------+---------+\n| ID | VOLUME | ISSUE | YEAR | SEASON  | PAGE | CATEGORY | MEMBER | INTERNATIONAL | REVIEW | TEXT                                      | SKETCH | FLAGGED |\n+----+--------+-------+------+---------+------+----------+--------+---------------+--------+-------------------------------------------+--------+---------+\n|  1 |      1 |     1 | 1986 | Mollit  |    1 | Pariatur |   1234 |               |        | Esse Lorem do nulla sunt mollit nulla in. |        |    ✔    |\n|  2 |      1 |     1 | 1986 | Eiusmod |    2 | Commodo  |  1234B |               |        | Magna officia anim dolore enim.           |        |    ✔    |\n+----+--------+-------+------+---------+------+----------+--------+---------------+--------+-------------------------------------------+--------+---------+\n\n+------------------------------------------------------+\n| Correspondence Matches:                              |\n+-----------+--------+----------+------------+---------+\n| REFERENCE | SENDER | RECEIVER | DATE       | LINK    |\n+-----------+--------+----------+------------+---------+\n| 123d5f    |     55 |     1234 | 1986-04-01 |    L1   |\n| b12cd3    |   1234 |       55 | 1986-05-16 | M123d5f |\n| 6beef9    |   1234 |      666 | 2021-03-15 |         |\n+-----------+--------+----------+------------+---------+\n",
		},
		{
			name:      "no listings, no correspondence",
			args:      []string{"42"},
			assertion: assert.NoError,
			want:      "\nNo LEX listings found.\n\nNo correspondences found.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd.NewSearchCmd()
			b := bytes.NewBufferString("")
			cmd.SetOut(b)
			cmd.SetArgs(tt.args)
			tt.assertion(t, cmd.Execute())
			out, err := io.ReadAll(b)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(out))
		})
	}
}

func TestSearchListings(t *testing.T) {
	m, dbFilePath, err := initDatastoreManager()
	assert.NoError(t, err)

	defer func() {
		m.Stop()
		assert.NoError(t, os.Remove(dbFilePath))
	}()

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
			want:    []lstg.Listing{},
			wantErr: false,
		},
		// TODO: checked against max_results in config
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.SearchListings(tt.args.member, m)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func initDatastoreManager() (*datastore.Manager, string, error) {
	currentTime := time.Now()
	filename := fmt.Sprintf("test_%d.db", currentTime.Unix())
	manager, err := datastore.New(filename)
	if err != nil {
		return nil, "", err
	}

	appFS := afero.NewMemMapFs()

	// create test files and directories
	err = appFS.MkdirAll("test", 0o755)
	if err != nil {
		return nil, "", err
	}

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

	return manager, filename, nil
}

func init() {
	viper.GetViper().Set("search.max_results", 10)
}
