// Application which greets you.
package main_test

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	main "github.com/asphaltbuffet/ogma"
)

func TestInitConfig(t *testing.T) {
	type args struct {
		cf string
	}
	tests := []struct {
		name      string
		args      args
		wantLevel log.Level
		assertion assert.ErrorAssertionFunc
	}{
		// {
		// 	name: "debug logging", // TODO: add afero to swap out config files for more test cases
		// 	args: args{
		// 		cf: ".ogma",
		// 	},
		// 	wantLevel: log.DebugLevel,
		// 	assertion: assert.NoError,
		// },
		{
			name: "warn logging",
			args: args{
				cf: ".ogma",
			},
			wantLevel: log.WarnLevel,
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main.InitConfig(tt.args.cf)
			got, err := log.ParseLevel(viper.GetString("logging.level"))
			tt.assertion(t, err)
			if err == nil {
				assert.Equal(t, tt.wantLevel, got)
			}
		})
	}
}
