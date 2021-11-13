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
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/asphaltbuffet/ogma/cmd"
	cmd2 "github.com/asphaltbuffet/ogma/pkg/cmd2"
	dsmocks "github.com/asphaltbuffet/ogma/pkg/datastore/mocks"
)

func TestRunSearchCmd(t *testing.T) {
	type args struct {
		c    *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cmd.RunSearchCmd(tt.args.c, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("RunSearchCmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name      string
		member    int
		want      []cmd2.Listing
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDatastore := &dsmocks.Reader{}
			mockDatastore.On("Data", mock.Anything).Return()

			got, err := cmd.Search(tt.member, mockDatastore)
			tt.assertion(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, got)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search() = %v, want %v", got, tt.want)
			}
		})
	}
}
