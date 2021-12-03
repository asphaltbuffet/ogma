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
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/ogma/cmd"
	"github.com/asphaltbuffet/ogma/pkg/datastore"
)

func TestNewImportCmd(t *testing.T) {
	got := cmd.NewImportCmd()

	assert.Equal(t, "import", got.Name())
	assert.Equal(t, "Bulk import records.", got.Short)
	assert.False(t, got.Runnable())
}

func Setup(t *testing.T) (*datastore.Manager, string, afero.Fs) {
	t.Helper()

	appFS := afero.NewOsFs()

	// create test files and directories
	assert.NoError(t, appFS.MkdirAll("test", 0o755))

	currentTime := time.Now()
	filename := fmt.Sprintf("test/test_%d.db", currentTime.Unix())
	manager, err := datastore.New(filename)
	assert.NoError(t, err)

	err = afero.WriteFile(appFS, "test/invalid.json", []byte(`{
		"listings": [
			{
				"volume": 2,
				"issue": 55,
				"year": 2021,
				"season": "Spring",
				"page": 1,
				"category": "Art & Photography",
				"member": 2989,
				"alt": "",
		]
		}`), 0o644)
	assert.NoError(t, err)

	err = afero.WriteFile(appFS, "test/listing.json", []byte(`{
		"listings": [
		{
		"volume": 2,
		"issue": 55,
		"year": 2021,
		"season": "Spring",
		"page": 1,
		"category": "Art & Photography",
		"member": 2989,
		"alt": "B",
		"international": false,
		"review": false,
		"text": "Fingerpainting exchange.",
		"art": false,
		"flag": false
		}
		]
		}`), 0o644)
	assert.NoError(t, err)

	err = afero.WriteFile(appFS, "test/listings.json", []byte(`{
		"listings": [
			{
				"volume": 1,
				"issue": 1,
				"year": 1986,
				"season": "Mollit",
				"page": 1,
				"category": "Pariatur",
				"member": 1234,
				"alt": "",
				"international": false,
				"review": false,
				"text": "Esse Lorem do nulla sunt mollit nulla in.",
				"art": false,
				"flag": true
			},
			{
				"volume": 1,
				"issue": 1,
				"year": 1986,
				"season": "Eiusmod",
				"page": 2,
				"category": "Commodo",
				"member": 1234,
				"alt": "B",
				"international": false,
				"review": false,
				"text": "Magna officia anim dolore enim.",
				"art": false,
				"flag": true
			},
			{
				"volume": 1,
				"issue": 1,
				"year": 1986,
				"season": "Id",
				"page": 3,
				"category": "Pariatur",
				"member": 5678,
				"alt": "",
				"international": false,
				"review": false,
				"text": "Velit cillum cillum ea officia nulla enim.",
				"art": false,
				"flag": true
			}
		]
		}`), 0o644)
	assert.NoError(t, err)

	err = afero.WriteFile(appFS, "test/mails.json", []byte(`{
		"mails": [
			{
				"reference": "123d5f",
				"sender": 55,
				"receiver": 1234,
				"date": "1986-04-01",
				"link": "L1"
			},
			{
				"reference": "b12cd3",
				"sender": 1234,
				"receiver": 55,
				"date": "1986-05-16",
				"link": "M123d5f"
			},
			{
				"reference": "6beef9",
				"sender": 1234,
				"receiver": 666,
				"date": "2021-03-15",
				"link": ""
			}
		]
		}`), 0o644)
	assert.NoError(t, err)

	return manager, filename, appFS
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
