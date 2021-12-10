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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/asphaltbuffet/ogma/cmd"
	"github.com/asphaltbuffet/ogma/mocks"
)

func TestNewImportMailCmd(t *testing.T) {
	got := cmd.NewImportMailCmd()

	assert.Equal(t, "mail", got.Name())
	assert.Equal(t, "Bulk import mail records.", got.Short)
	assert.True(t, got.Runnable())
}

func TestRunImportMailCmd(t *testing.T) {
	m, dbFilePath, appFS := Setup(t)
	m.Stop()

	defer func() {
		assert.NoError(t, appFS.RemoveAll("test/"))
	}()

	viper.Set("datastore.filename", dbFilePath)

	tests := []struct {
		name      string
		args      []string
		assertion assert.ErrorAssertionFunc
		want      string
	}{
		{
			name:      "mail import",
			args:      []string{"test/mails.json"},
			assertion: assert.NoError,
			want:      "Imported 3/3 mail records.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd.NewImportMailCmd()
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

func TestImportMail(t *testing.T) {
	m, _, appFS := Setup(t)
	m.Stop()

	defer func() {
		assert.NoError(t, appFS.RemoveAll("test/"))
	}()

	tests := []struct {
		name      string
		filepath  string
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "single entry",
			filepath:  "test/mails.json",
			want:      "Imported 3/3 mail records.",
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile, err := appFS.Open(tt.filepath)
			assert.NoError(t, err)

			mockDatastore := &mocks.Writer{}
			mockDatastore.On("Save", mock.Anything).Return(nil)

			got, err := cmd.ImportMail(testFile, mockDatastore)
			tt.assertion(t, err)

			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestParseMail(t *testing.T) {
	m, _, appFS := Setup(t)
	m.Stop()

	defer func() {
		require.NoError(t, appFS.RemoveAll("test/"))
	}()

	type args struct {
		fp string
	}
	tests := []struct {
		name      string
		args      args
		want      []cmd.Mail
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "multiple mail entry",
			args: args{
				fp: "test/mails.json",
			},
			want: []cmd.Mail{
				{Ref: "123d5f", Sender: 55, Receiver: 1234, Date: "1986-04-01", Link: "L1"},
				{Ref: "b12cd3", Sender: 1234, Receiver: 55, Date: "1986-05-16", Link: "M123d5f"},
				{Ref: "6beef9", Sender: 1234, Receiver: 666, Date: "2021-03-15", Link: ""},
			},
			assertion: assert.NoError,
		},
		{
			name: "invalid json",
			args: args{
				fp: "test/invalid.json",
			},
			want:      []cmd.Mail{},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile, err := appFS.Open(tt.args.fp)
			require.NoError(t, err)

			got, err := cmd.ParseMail(testFile)
			tt.assertion(t, err)

			if err == nil {
				assert.Truef(t, assert.ObjectsAreEqual(tt.want, got), "%s: expected %+v, got %+v", tt.name, tt.want, got)
			}
		})
	}
}

func TestUniqueMail(t *testing.T) {
	type args struct {
		mm []cmd.Mail
	}
	tests := []struct {
		name string
		args args
		want []cmd.Mail
	}{
		{
			name: "empty",
			args: args{
				mm: []cmd.Mail{},
			},
			want: []cmd.Mail{},
		},
		{
			name: "no duplicates",
			args: args{
				mm: []cmd.Mail{
					{ID: 0, Ref: "", Sender: 0, Receiver: 0, Date: "", Link: ""},
				},
			},
			want: []cmd.Mail{
				{ID: 0, Ref: "", Sender: 0, Receiver: 0, Date: "", Link: ""},
			},
		},
		{
			name: "only duplicates",
			args: args{
				mm: []cmd.Mail{
					{ID: 0, Ref: "", Sender: 0, Receiver: 0, Date: "", Link: ""},
				},
			},
			want: []cmd.Mail{
				{ID: 0, Ref: "", Sender: 0, Receiver: 0, Date: "", Link: ""},
			},
		},
		{
			name: "duplicates with unique",
			args: args{
				mm: []cmd.Mail{
					{ID: 0, Ref: "", Sender: 0, Receiver: 0, Date: "", Link: ""},
				},
			},
			want: []cmd.Mail{
				{ID: 0, Ref: "", Sender: 0, Receiver: 0, Date: "", Link: ""},
			},
		},
		{
			name: "multiple duplicates with unique",
			args: args{
				mm: []cmd.Mail{
					{ID: 0, Ref: "", Sender: 0, Receiver: 0, Date: "", Link: ""},
				},
			},
			want: []cmd.Mail{
				{ID: 0, Ref: "", Sender: 0, Receiver: 0, Date: "", Link: ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmd.UniqueMails(tt.args.mm); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UniqueMails() = %v, want %v", got, tt.want)
			}
		})
	}
}
