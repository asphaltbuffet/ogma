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
	"io/ioutil"
	"os"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/ogma/cmd"
)

func TestNewMailCmd(t *testing.T) {
	got := cmd.NewMailCmd()

	assert.Equal(t, "mail", got.Name())
	assert.Equal(t, "Tracks letters sent to/from penpals", got.Short)
	assert.True(t, got.Runnable())
}

func TestRunMailCmd(t *testing.T) {
	currentTime := time.Now()
	filename := fmt.Sprintf("test_%d.db", currentTime.Unix())

	defer func() {
		assert.NoError(t, os.Remove(filename))
	}()

	viper.Set("datastore.filename", filename)
	tests := []struct {
		name      string
		args      []string
		assertion assert.ErrorAssertionFunc
		want      string
	}{
		{
			name:      "valid",
			args:      []string{"-d2021-11-15", "-s1234", "-r5678"},
			assertion: assert.NoError,
			want:      "Added mail. Reference: f2165e\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd.NewMailCmd()
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
				Date:     time.Date(2021, time.November, 15, 0, 0, 0, 0, time.Local),
			},
			length: 6,
			want:   "f2165e",
		},
		{
			name: "6 char hash - 2", // try with different values to ensure we're getting variation
			mail: cmd.Mail{
				Sender:   123,
				Receiver: 45678,
				Date:     time.Date(2021, time.November, 15, 0, 0, 0, 0, time.Local),
			},
			length: 6,
			want:   "650e0a",
		},
		{
			name: "0 char hash",
			mail: cmd.Mail{
				Sender:   1234,
				Receiver: 5678,
				Date:     time.Date(2021, time.November, 15, 0, 0, 0, 0, time.Local),
			},
			length: 0,
			want:   "",
		},
		{
			name: "full hash",
			mail: cmd.Mail{
				Sender:   1234,
				Receiver: 5678,
				Date:     time.Date(2021, time.November, 15, 0, 0, 0, 0, time.Local),
			},
			length: 32,
			want:   "28bf0b58528e41181e13d0f789f2165e",
		},
		{
			name: "overbound hash",
			mail: cmd.Mail{
				Sender:   1234,
				Receiver: 5678,
				Date:     time.Date(2021, time.November, 15, 0, 0, 0, 0, time.Local),
			},
			length: 33,
			want:   "28bf0b58528e41181e13d0f789f2165e",
		},
		{
			name: "underbound hash",
			mail: cmd.Mail{
				Sender:   1234,
				Receiver: 5678,
				Date:     time.Date(2021, time.November, 15, 0, 0, 0, 0, time.Local),
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
		want      time.Time
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "good date",
			date:      "2021-11-15",
			want:      time.Date(2021, time.November, 15, 0, 0, 0, 0, time.Local),
			assertion: assert.NoError,
		},
		{
			name:      "bad date",
			date:      "Nov 15 2021",
			want:      time.Date(2021, time.November, 15, 0, 0, 0, 0, time.Local),
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmd.ValidateDate(tt.date)
			tt.assertion(t, err)
			if err == nil {
				assert.Truef(t, got.Equal(tt.want), "Times are not equal: got = %v, want %v", got, tt.want)
			}
		})
	}
}

func init() {
	log.SetOutput(ioutil.Discard)
}
