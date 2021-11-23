// Application which greets you.
package main_test

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	main "github.com/asphaltbuffet/ogma"
)

func TestInitConfig(t *testing.T) {
	var err error

	appFS := afero.NewOsFs()
	viper.AddConfigPath("test/")
	defer func() {
		err = appFS.RemoveAll("test/")
		assert.NoError(t, err)
	}()

	// create test files and directories
	err = appFS.MkdirAll("test", 0o755)
	assert.NoError(t, err)

	err = afero.WriteFile(appFS, "test/.debugConfig", []byte("logging:\n"+
		"    level: \"debug\"\n"+
		"search:\n"+
		"    max_results: 10\n"+
		"datastore:\n"+
		"    filename: \"ogma.db\"\n"+
		"defaults:\n"+
		"    issue: 56\n"+
		"    max_column: 40\n"), 0o644)
	assert.NoError(t, err)

	err = afero.WriteFile(appFS, "test/.errorConfig", []byte("logging:\n"+
		"    level: \"error\"\n"+
		"search:\n"+
		"    max_results: 10\n"+
		"datastore:\n"+
		"    filename: \"ogma.db\"\n"+
		"defaults:\n"+
		"    issue: 56\n"+
		"    max_column: 40\n"), 0o644)
	assert.NoError(t, err)
	log.Warn(afero.Exists(appFS, ".errorConfig"))

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
			name: "filename doesn't exist",
			args: args{
				cf: "foo",
			},
			wantLevel: log.WarnLevel,
			assertion: assert.Error,
		},
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
			main.InitConfig(appFS, tt.args.cf)
			got, err := log.ParseLevel(viper.GetString("logging.level"))
			tt.assertion(t, err)
			if err == nil {
				assert.Equal(t, tt.wantLevel, got)
			}
		})
	}
}
