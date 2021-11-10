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
	cmd2 "github.com/asphaltbuffet/ogma/pkg/cmd2"
)

var verbose bool

func TestListingsSearchCmd(t *testing.T) {
	var buf bytes.Buffer

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	assert.NoError(t, err)

	tt := map[string]struct {
		args     []string
		function func() func(cmd *cobra.Command, args []string) error
		output   string
	}{
		"all default flags": {
			args: []string{"listings", "search"},
			function: func() func(c *cobra.Command, args []string) error {
				return func(c *cobra.Command, args []string) error {
					c.SetOut(&buf)
					err := cmd2.RunSearchListings(c)
					assert.NoError(t, err)
					return nil
				}
			},
			output: "Found 0 listings.\n",
		},
		"no results by year (short)": {
			args: []string{"listings", "search", "-y1"},
			function: func() func(c *cobra.Command, args []string) error {
				return func(c *cobra.Command, args []string) error {
					c.SetOut(&buf)
					err := cmd2.RunSearchListings(c)
					assert.NoError(t, err)
					return nil
				}
			},
			output: "Found 0 listings.\n",
		},
		"no results by year (long)": {
			args: []string{"listings", "search", "--year=1"},
			function: func() func(c *cobra.Command, args []string) error {
				return func(c *cobra.Command, args []string) error {
					c.SetOut(&buf)
					err := cmd2.RunSearchListings(c)
					assert.NoError(t, err)
					return nil
				}
			},
			output: "Found 0 listings.\n",
		},
		"valid search by year": {
			args: []string{"listings", "search", "-y2021"},
			function: func() func(c *cobra.Command, args []string) error {
				return func(c *cobra.Command, args []string) error {
					c.SetOut(&buf)
					err := cmd2.RunSearchListings(c)
					assert.NoError(t, err)
					return nil
				}
			},
			output: "Found 3 listings.\n",
		},
		"multi-flag search": {
			args: []string{"listings", "search", "-y2021", "-i56"},
			function: func() func(c *cobra.Command, args []string) error {
				return func(c *cobra.Command, args []string) error {
					c.SetOut(&buf)
					err := cmd2.RunSearchListings(c)
					assert.NoError(t, err)
					return nil
				}
			},
			output: "Found 10 listings.\n",
		},
		"search results excede max limitation": {
			args: []string{"listings", "search", "-y2021", "-i56", "-m1000"},
			function: func() func(c *cobra.Command, args []string) error {
				return func(c *cobra.Command, args []string) error {
					c.SetOut(&buf)
					err := cmd2.RunSearchListings(c)
					assert.Error(t, err)
					return nil
				}
			},
			output: "Found 12 listings.\n",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			ogmaCmd := &cobra.Command{Use: "ogma"}

			listingsCmd := &cobra.Command{
				Use: "listings",
				RunE: func(c *cobra.Command, args []string) error {
					return nil
				},
			}

			searchListingsCmd := &cobra.Command{
				Use:  "search",
				Args: cobra.NoArgs,
				RunE: tc.function(),
			}

			searchListingsCmd.Flags().IntP("year", "y", -1, "Search listings by LEX Issue year.")
			searchListingsCmd.Flags().IntP("issue", "i", -1, "Search listings by LEX Issue Number.")
			searchListingsCmd.Flags().IntP("member", "m", -1, "Search listings by member number.")
			listingsCmd.AddCommand(searchListingsCmd)
			ogmaCmd.AddCommand(listingsCmd)

			c, out, err := ExecuteCommandC(t, ogmaCmd, tc.args...)
			if out != "" {
				t.Errorf("Unexpected output: %v", out)
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.output, buf.String())
			if c.Name() != "search" {
				t.Errorf(`invalid command returned from ExecuteC: expected "search"', got: %q`, c.Name())
			}
			buf.Reset()
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
