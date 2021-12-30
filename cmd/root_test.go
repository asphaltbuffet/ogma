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

// Package cmd contains all CLI commands used by the application.
package cmd_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/asphaltbuffet/ogma/cmd"
)

func TestGenerateInfo(t *testing.T) {
	m, dbfilename, fs := setup(t)
	m.Stop()

	defer func() {
		require.NoError(t, fs.RemoveAll("test/"))
	}()

	tests := []struct {
		name      string
		args      []string
		assertion assert.ErrorAssertionFunc
		want      string
	}{
		{
			name:      "valid - empty db",
			args:      []string{"-p=false"},
			assertion: assert.NoError,
			want:      "+--------------+\n| Data Records |\n| :            |\n+----------+---+\n| Mail     | 0 |\n| Listings | 0 |\n+----------+---+\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("datastore.filename", dbfilename)

			c := cmd.GetRootCmd()

			b := bytes.NewBufferString("")
			c.SetOut(b)
			c.SetArgs(tt.args)

			tt.assertion(t, c.Execute())

			out, err := io.ReadAll(b)
			require.NoError(t, err)

			assert.Contains(t, string(out), tt.want)
		})
	}
}

func TestInitConfig(t *testing.T) {
	var err error

	appFS := afero.NewOsFs()
	defer func() {
		err = appFS.RemoveAll("test/")
		require.NoError(t, err)
	}()

	// create test files and directories
	err = appFS.MkdirAll("test", 0o755)
	require.NoError(t, err)

	err = afero.WriteFile(appFS, "test/.debugConfig", []byte("logging:\n"+
		"    level: \"debug\"\n"+
		"search:\n"+
		"    max_results: 10\n"+
		"datastore:\n"+
		"    filename: \"ogma.db\"\n"+
		"defaults:\n"+
		"    issue: 56\n"+
		"    max_column: 40\n"), 0o644)
	require.NoError(t, err)

	err = afero.WriteFile(appFS, "test/.errorConfig", []byte("logging:\n"+
		"    level: \"error\"\n"+
		"search:\n"+
		"    max_results: 10\n"+
		"datastore:\n"+
		"    filename: \"ogma.db\"\n"+
		"defaults:\n"+
		"    issue: 56\n"+
		"    max_column: 40\n"), 0o644)
	require.NoError(t, err)

	type args struct {
		cf string
	}
	tests := []struct {
		name      string
		args      args
		wantLevel log.Level
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "debug logging",
			args: args{
				cf: ".debugConfig",
			},
			wantLevel: log.DebugLevel,
			assertion: assert.NoError,
		},
		{
			name: "error logging",
			args: args{
				cf: ".errorConfig",
			},
			wantLevel: log.ErrorLevel,
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.FileExists(t, fmt.Sprintf("test/%s", tt.args.cf), "config file not found")

			viper.AddConfigPath("test/")
			cmd.InitConfig(appFS, tt.args.cf)

			got, err := log.ParseLevel(viper.GetString("logging.level"))
			tt.assertion(t, err)
			if err == nil {
				assert.Equal(t, tt.wantLevel, got)
			}
		})
	}
}
