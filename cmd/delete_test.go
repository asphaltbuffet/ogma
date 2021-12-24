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
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/asphaltbuffet/ogma/cmd"
)

func TestNewDeleteCmd(t *testing.T) {
	got := cmd.NewDeleteCmd()

	assert.Equal(t, "delete", got.Name())
	assert.True(t, got.Runnable())
}

func TestRunDeleteCmd(t *testing.T) {
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
			name:      "successful deletion",
			args:      []string{"-a"},
			datastore: dsFile,
			assertion: assert.NoError,
			want:      "all data records have been removed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("datastore.filename", tt.datastore)

			cmd := cmd.NewDeleteCmd()
			b := bytes.NewBufferString("")
			cmd.SetOut(b)
			cmd.SetErr(b)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			tt.assertion(t, err)

			// assert.Equal(t, tt.want, b.String()[:len(tt.want)], "unexpected output")
			assert.Contains(t, b.String(), tt.want)
		})
	}
}
