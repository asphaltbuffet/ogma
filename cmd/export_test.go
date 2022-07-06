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

func TestNewExportCmd(t *testing.T) {
	got := cmd.NewExportCmd()

	assert.Equal(t, "export", got.Name())
	assert.True(t, got.Runnable())
}

func TestRunExportCmd(t *testing.T) {
	m, dsFile := initDatastoreManager(t)
	m.Stop()

	defer func() {
		require.NoError(t, os.RemoveAll("test/"))
	}()

	validExportFile := "test/export.json"
	validExportArg := "-o=" + validExportFile
	invalidExportFile := "test/test/export.txt"
	invalidExportArg := "-o=" + invalidExportFile
	successReturn := "successfully exported"
	errorReturn := "error exporting"

	tests := []struct {
		name      string
		args      []string
		datastore string
		export    string
		assertion assert.ErrorAssertionFunc
		want      string
	}{
		{
			name:      "mail export",
			args:      []string{"-r=mail", validExportArg},
			datastore: dsFile,
			export:    "test/export.json",
			assertion: assert.NoError,
			want:      successReturn,
		},
		{
			name:      "listing export",
			args:      []string{"-r=listing", validExportArg},
			datastore: dsFile,
			export:    "test/export.json",
			assertion: assert.NoError,
			want:      successReturn,
		},
		{
			name:      "all export",
			args:      []string{"-r=all", validExportArg},
			datastore: dsFile,
			export:    "test/export.json",
			assertion: assert.NoError,
			want:      successReturn,
		},
		{
			name:      "mail export - invalid path for export file",
			args:      []string{"-r=mail", invalidExportArg},
			datastore: dsFile,
			export:    "test/test/export.json",
			assertion: assert.Error,
			want:      errorReturn,
		},
		{
			name:      "listing export - invalid path for export file",
			args:      []string{"-r=listing", invalidExportArg},
			datastore: dsFile,
			export:    "test/test/export.json",
			assertion: assert.Error,
			want:      errorReturn,
		},
		{
			name:      "all export - invalid path for export file",
			args:      []string{"-r=all", invalidExportArg},
			datastore: dsFile,
			export:    "test/test/export.json",
			assertion: assert.Error,
			want:      errorReturn,
		},
		{
			name:      "invalid record type",
			args:      []string{"-r=foo", validExportArg},
			datastore: dsFile,
			export:    "test/test/export.json",
			assertion: assert.Error,
			want:      "invalid option:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("datastore.filename", tt.datastore)

			cmd := cmd.NewExportCmd()
			b := bytes.NewBufferString("")
			cmd.SetOut(b)
			cmd.SetErr(b)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			require.NoError(t, err)

			_, err = os.Stat(tt.export)
			tt.assertion(t, err)

			assert.Contains(t, b.String(), tt.want)
		})
	}
}
