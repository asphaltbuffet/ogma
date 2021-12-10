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
	"io/ioutil"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/asphaltbuffet/ogma/cmd"
)

func TestNewMailCmd(t *testing.T) {
	got := cmd.NewMailCmd()

	assert.Equal(t, "mail", got.Name())
	assert.Equal(t, "Tracks letters sent to/from penpals", got.Short)
	assert.True(t, got.Runnable())
}

func TestRunMailCmd(t *testing.T) {
	m, dbfilename, fs := Setup(t)
	m.Stop()

	defer func() {
		// assert.NoError(t, os.Remove(dbfilename))
		require.NoError(t, fs.RemoveAll("test/"))
	}()

	viper.Set("datastore.filename", dbfilename)
	tests := []struct {
		name      string
		args      []string
		config    string
		assertion assert.ErrorAssertionFunc
		want      string
	}{
		{
			name:      "valid",
			args:      []string{"-d2021-11-15", "-s1234", "-r5678"},
			config:    "/test/.tconfig",
			assertion: assert.NoError,
			want:      "Added mail. Reference: f2165e\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// cmd.NewRootCmd()
			c := cmd.NewMailCmd()
			b := bytes.NewBufferString("")
			cmd.InitConfig(fs, tt.config)
			c.SetOut(b)
			c.SetArgs(tt.args)
			tt.assertion(t, c.Execute())
			out, err := io.ReadAll(b)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(out))
		})
	}
}

func TestMailHash(t *testing.T) {
	tests := []struct {
		name   string
		mail   cmd.Mail
		length int
		want   string
	}{
		{
			name: "6 char hash",
			mail: cmd.Mail{
				Sender:   1234,
				Receiver: 5678,
				Date:     "2021-11-15",
			},
			length: 6,
			want:   "f2165e",
		},
		{
			name: "6 char hash - 2", // try with different values to ensure we're getting variation
			mail: cmd.Mail{
				Sender:   123,
				Receiver: 45678,
				Date:     "2021-11-15",
			},
			length: 6,
			want:   "650e0a",
		},
		{
			name: "0 char hash",
			mail: cmd.Mail{
				Sender:   1234,
				Receiver: 5678,
				Date:     "2021-11-15",
			},
			length: 0,
			want:   "",
		},
		{
			name: "full hash",
			mail: cmd.Mail{
				Sender:   1234,
				Receiver: 5678,
				Date:     "2021-11-15",
			},
			length: 32,
			want:   "28bf0b58528e41181e13d0f789f2165e",
		},
		{
			name: "overbound hash",
			mail: cmd.Mail{
				Sender:   1234,
				Receiver: 5678,
				Date:     "2021-11-15",
			},
			length: 33,
			want:   "28bf0b58528e41181e13d0f789f2165e",
		},
		{
			name: "underbound hash",
			mail: cmd.Mail{
				Sender:   1234,
				Receiver: 5678,
				Date:     "2021-11-15",
			},
			length: -1,
			want:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmd.MailHash(tt.mail, tt.length); got != tt.want {
				t.Errorf("MailHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateDate(t *testing.T) {
	tests := []struct {
		name      string
		date      string
		want      string
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "good date",
			date:      "2021-11-15",
			want:      "2021-11-15",
			assertion: assert.NoError,
		},
		{
			name:      "bad date",
			date:      "Nov 15 2021",
			want:      "2021-11-15",
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.ValidateDate(tt.date)
			tt.assertion(t, err)
			if err == nil {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func init() {
	log.SetOutput(ioutil.Discard)
}
