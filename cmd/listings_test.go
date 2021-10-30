package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestListingsSearchCmd(t *testing.T) { //nolint:funlen // ignore this for now 2021/10/29 BL
	var buf bytes.Buffer

	// Search config in application directory with name ".ogma" (without extension).
	viper.AddConfigPath("../")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".ogma")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	assert.NoError(t, err)

	tests := []struct {
		args     []string
		function func() func(cmd *cobra.Command, args []string) error
		output   string
	}{
		// no flags returns all results (will be an error if db has more than MAX_RESULTS listings)
		{
			args: []string{"listings", "search", "unknown"},
			function: func() func(cmd *cobra.Command, args []string) error {
				return func(cmd *cobra.Command, args []string) error {
					cmd.SetOut(&buf)
					err := RunSearchListings(cmd)
					assert.NoError(t, err)
					return nil
				}
			},
			output: "Found 0 listings.\n",
		},
		// // search for year less than min
		// {
		// 	args: []string{"-y1"},
		// 	err:  nil,
		// 	out:  "Found 0 listings.",
		// },
		// {
		// 	args: []string{"--year=1"},
		// 	err:  nil,
		// 	out:  "Found 0 listings.",
		// },
		// // search for year greater than max
		// {
		// 	args: []string{"-y9999"},
		// 	err:  nil,
		// 	out:  "Found 0 listings.",
		// },
		// {
		// 	args: []string{"--year=9999"},
		// 	err:  nil,
		// 	out:  "Found 0 listings.",
		// },
	}

	for _, testcase := range tests {
		ogmaCmd := &cobra.Command{Use: "ogma"}

		listingsCmd = &cobra.Command{
			Use: "listings",
			RunE: func(cmd *cobra.Command, args []string) error {
				return nil
			},
		}

		searchListingsCmd := &cobra.Command{
			Use:  "search",
			RunE: testcase.function(),
		}

		searchListingsCmd.Flags().IntP("year", "y", -1, "Search listings by LEX Issue year.")
		searchListingsCmd.Flags().IntP("issue", "i", -1, "Search listings by LEX Issue Number.")
		searchListingsCmd.Flags().IntP("member", "m", -1, "Search listings by member number.")
		listingsCmd.AddCommand(searchListingsCmd)
		ogmaCmd.AddCommand(listingsCmd)

		c, out, err := ExecuteCommandC(t, ogmaCmd, testcase.args...)
		if out != "" {
			t.Errorf("Unexpected output: %v", out)
		}
		assert.NoError(t, err)
		assert.Equal(t, testcase.output, buf.String())
		if c.Name() != "search" {
			t.Errorf(`invalid command returned from ExecuteC: expected "search"', got: %q`, c.Name())
		}
		buf.Reset()
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
