package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	tt := []struct {
		args []string
		err  error
		out  string
	}{
		{
			args: nil,
			err:  errors.New("not ok"),
		},
		{
			args: []string{"-i"},
			err:  nil,
			out:  "ok",
		},
		{
			args: []string{"--info"},
			err:  nil,
			out:  "ok",
		},
	}

	root := &cobra.Command{Use: "root", RunE: RootCmdRunE}
	RootCmdFlags(root)

	for _, tc := range tt {
		out, err := execute(t, root, tc.args...)

		assert.Equal(t, tc.err, err)

		if tc.err == nil {
			assert.Equal(t, tc.out, out)
		}
	}
}

func execute(t *testing.T, c *cobra.Command, args ...string) (string, error) {
	t.Helper()

	buf := new(bytes.Buffer)
	c.SetOut(buf)
	c.SetErr(buf)
	c.SetArgs(args)

	err := c.Execute()
	return strings.TrimSpace(buf.String()), err
}
