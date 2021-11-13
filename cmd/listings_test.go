package cmd_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/ogma/cmd"
)

var verbose bool

func ExecuteCommandC(t *testing.T, root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	t.Helper()

	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func TestRunImportListingsCmd(t *testing.T) {
	var buf bytes.Buffer

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	assert.NoError(t, err)

	// Put the test_db in same place as config file for testing
	fp := filepath.Dir(viper.ConfigFileUsed())
	testDsfile := fp + "/test_db.db"

	// Change datastore for testing
	viper.Set("datastore.filename", testDsfile)
	defer func() {
		err = os.Remove(testDsfile)
		assert.NoError(t, err)
	}()
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		// // this testing doesn't allow verifying cobra behavior so far
		// {
		// 	name:    "no args",
		// 	args:    []string{"listings", "import"},
		// 	want:    "",
		// 	wantErr: true,
		// },
		{
			name:    "missing file",
			args:    []string{"listings", "import", fp + "/bad_file.json"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "good import",
			args:    []string{"listings", "import", fp + "/importSingle_test.json"},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ogmaCmd := &cobra.Command{Use: "ogma"}
			listingsCmd := &cobra.Command{
				Use: "listings",
				RunE: func(c *cobra.Command, args []string) error {
					return nil
				},
			}

			importListingsCmd := &cobra.Command{
				Use: "import",
				RunE: func(c *cobra.Command, args []string) error {
					ogmaCmd.SetOut(&buf)
					err := cmd.RunImportListingsCmd(c, args)
					if (err != nil) != tt.wantErr {
						t.Errorf("RunImportListingsCmd() error = %v, wantErr %v", err, tt.wantErr)
					}
					return nil
				},
			}

			importListingsCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print imported listings to stdout.")
			listingsCmd.AddCommand(importListingsCmd)
			listingsCmd.AddCommand(importListingsCmd)
			ogmaCmd.AddCommand(listingsCmd)

			c, out, err := ExecuteCommandC(t, ogmaCmd, tt.args...)
			assert.Emptyf(t, out, "Unexpected output: %v", out)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, buf.String())
			assert.Equalf(t, "import", c.Name(), `Invalid command returned from ExecuteC: expected "import", got: %q`, c.Name())
		})
	}
}

func init() {
	log.SetOutput(ioutil.Discard)

	// Search config in application directory with name ".ogma" (without extension).
	viper.AddConfigPath("../")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".ogma")

	viper.AutomaticEnv() // read in environment variables that match
}
