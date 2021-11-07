package cmd_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	cmd2 "github.com/asphaltbuffet/ogma/pkg/cmd2"
)

var (
	volume        int
	lex           int
	year          int
	season        string
	page          int
	category      string
	member        int
	international bool
	review        bool
	text          string
	art           bool
	flag          bool
)

func TestListingsSearchCmd(t *testing.T) {
	var buf bytes.Buffer

	// Search config in application directory with name ".ogma" (without extension).
	viper.AddConfigPath("../")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".ogma")

	viper.AutomaticEnv() // read in environment variables that match

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

func TestListingsAddCmd(t *testing.T) {
	var buf bytes.Buffer

	// Search config in application directory with name ".ogma" (without extension).
	viper.AddConfigPath("../")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".ogma")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	assert.NoError(t, err)

	// Change datastore for testing
	viper.Set("datastore.filename", "test_db.db")
	defer func() {
		err = os.Remove("test_db.db")
		assert.NoError(t, err)
	}()

	tt := map[string]struct {
		args     []string
		function func() func(cmd *cobra.Command, args []string) error
		output   string
	}{
		"no flags": {
			args: []string{
				"listings", "add",
			},
			function: func() func(c *cobra.Command, args []string) error {
				return func(c *cobra.Command, args []string) error {
					c.SetOut(&buf)
					assert.NoError(t, err)

					out, err := cmd2.RunAddListing([]cmd2.Listing{ //nolint:govet // fine for testing right now?
						{
							Volume:              volume,
							IssueNumber:         lex,
							Year:                year,
							Season:              season,
							PageNumber:          page,
							IndexedCategory:     category,
							IndexedMemberNumber: member,
							MemberExtension:     "",
							IsInternational:     international,
							IsReview:            review,
							ListingText:         text,
							IsArt:               art,
							IsFlagged:           flag,
						},
					})
					if err == nil {
						c.Println(out)
					}

					return err
				}
			},
			output: "+--------+-------+------+--------+------+----------+--------+---------------+--------+------+--------+---------+\n" +
				"| VOLUME | ISSUE | YEAR | SEASON | PAGE | CATEGORY | MEMBER | INTERNATIONAL | REVIEW | TEXT | SKETCH | FLAGGED |\n" +
				"+--------+-------+------+--------+------+----------+--------+---------------+--------+------+--------+---------+\n" +
				"|     -1 |    56 | 2021 |        |   -1 |          |     -1 | false         | false  |      | false  | false   |\n" +
				"+--------+-------+------+--------+------+----------+--------+---------------+--------+------+--------+---------+\n",
		},
		"required flags": {
			args: []string{
				"listings", "add", "-v2", "-l40", "-y2021", "-sAutumn", "-p2", "-cCrafts", "-m12345", "-t\"Some text goes here.\"",
			},
			function: func() func(c *cobra.Command, args []string) error {
				return func(c *cobra.Command, args []string) error {
					c.SetOut(&buf)
					assert.NoError(t, err)

					out, err := cmd2.RunAddListing([]cmd2.Listing{ //nolint:govet // fine for testing right now?
						{
							Volume:              volume,
							IssueNumber:         lex,
							Year:                year,
							Season:              season,
							PageNumber:          page,
							IndexedCategory:     category,
							IndexedMemberNumber: member,
							MemberExtension:     "",
							IsInternational:     international,
							IsReview:            review,
							ListingText:         text,
							IsArt:               art,
							IsFlagged:           flag,
						},
					})
					if err == nil {
						c.Println(out)
					}

					return err
				}
			},
			output: "+--------+-------+------+--------+------+----------+--------+---------------+--------+------------------------+--------+---------+\n" +
				"| VOLUME | ISSUE | YEAR | SEASON | PAGE | CATEGORY | MEMBER | INTERNATIONAL | REVIEW | TEXT                   | SKETCH | FLAGGED |\n" +
				"+--------+-------+------+--------+------+----------+--------+---------------+--------+------------------------+--------+---------+\n" +
				"|      2 |    40 | 2021 | Autumn |    2 | Crafts   |  12345 | false         | false  | \"Some text goes here.\" | false  | false   |\n" +
				"+--------+-------+------+--------+------+----------+--------+---------------+--------+------------------------+--------+---------+\n",
		},
		"all flags": {
			args: []string{
				"listings", "add", "-v9", "-l999", "-y9999", "-sasdfb", "-p9", "-c\"some category\"", "-m9876", "-t\"Some kind of text goes here.\"", "-i", "-r", "-a", "-f",
			},
			function: func() func(c *cobra.Command, args []string) error {
				return func(c *cobra.Command, args []string) error {
					c.SetOut(&buf)
					assert.NoError(t, err)

					ll := []cmd2.Listing{
						{
							Volume:              volume,
							IssueNumber:         lex,
							Year:                year,
							Season:              season,
							PageNumber:          page,
							IndexedCategory:     category,
							IndexedMemberNumber: member,
							MemberExtension:     "",
							IsInternational:     international,
							IsReview:            review,
							ListingText:         text,
							IsArt:               art,
							IsFlagged:           flag,
						},
					}
					out, err := cmd2.RunAddListing(ll)
					if err == nil {
						c.Println(out)
					}

					return err
				}
			},
			output: "+--------+-------+------+--------+------+-----------------+--------+---------------+--------+--------------------------------+--------+---------+\n" +
				"| VOLUME | ISSUE | YEAR | SEASON | PAGE | CATEGORY        | MEMBER | INTERNATIONAL | REVIEW | TEXT                           | SKETCH | FLAGGED |\n" +
				"+--------+-------+------+--------+------+-----------------+--------+---------------+--------+--------------------------------+--------+---------+\n" +
				"|      9 |   999 | 9999 | asdfb  |    9 | \"some category\" |   9876 | true          | true   | \"Some kind of text goes here.\" | true   | true    |\n" +
				"+--------+-------+------+--------+------+-----------------+--------+---------------+--------+--------------------------------+--------+---------+\n",
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

			addListingCmd := &cobra.Command{
				Use:  "add",
				Args: cobra.NoArgs,
				RunE: tc.function(),
			}

			addListingCmd.Flags().IntVarP(&volume, "volume", "v", -1, "Volume containing listing entry.")
			addListingCmd.Flags().IntVarP(&lex, "lex", "l", viper.GetInt("defaults.issue"), "LEX issue containing listing entry.")
			addListingCmd.Flags().IntVarP(&year, "year", "y", time.Now().Year(), "Year of listing entry..")
			addListingCmd.Flags().StringVarP(&season, "season", "s", "", "Season of listing entry.")
			addListingCmd.Flags().IntVarP(&page, "page", "p", -1, "Page number of listing entry.")
			addListingCmd.Flags().StringVarP(&category, "category", "c", "", "Category of listing entry.")
			addListingCmd.Flags().IntVarP(&member, "member", "m", -1, "Member number of listing entry.")
			addListingCmd.Flags().BoolVarP(&international, "international", "i", false, "Is international postage required?")
			addListingCmd.Flags().BoolVarP(&review, "review", "r", false, "Is this a book review listing entry?")
			addListingCmd.Flags().StringVarP(&text, "text", "t", "", "Text of listing entry.")
			addListingCmd.Flags().BoolVarP(&art, "art", "a", false, "Is this a sketch listing entry?")
			addListingCmd.Flags().BoolVarP(&flag, "flag", "f", false, "Has this listing entry been flagged?")
			listingsCmd.AddCommand(addListingCmd)
			ogmaCmd.AddCommand(listingsCmd)

			c, out, err := ExecuteCommandC(t, ogmaCmd, tc.args...)
			if out != "" {
				t.Errorf("Unexpected output: %v", out)
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.output, buf.String())
			if c.Name() != "add" {
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

func init() {
	log.SetOutput(ioutil.Discard)
}
