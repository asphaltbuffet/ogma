package datastore_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/asdine/storm/v3/q"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
)

type testEntry struct {
	ID    int `storm:"id,increment"`
	Key   int `storm:"index"`
	Value string
}

func TestManagerNew(t *testing.T) {
	manager, dbFilePath := initDatastoreManager(t)

	defer func() {
		manager.Stop()
		err := os.Remove(dbFilePath)
		require.NoError(t, err)
	}()

	_, err := os.Stat(dbFilePath)
	assert.NotErrorIs(t, err, os.ErrNotExist)
}

func TestManagerOpen(t *testing.T) {
	manager, dbFilePath := initDatastoreManager(t)

	manager.Stop()
	defer func() {
		err := os.Remove(dbFilePath)
		require.NoError(t, err)
	}()

	m, err := datastore.Open(dbFilePath)
	assert.NoError(t, err, "should be able to open datastore that was just created")

	m.Stop()

	_, err = datastore.Open("foo.db")
	assert.Error(t, err, "should fail to open datastore that doesn't exist")

	manager.Stop()
}

func TestManagerNewFail(t *testing.T) {
	manager, dbFilePath := initDatastoreManager(t)
	defer func() {
		manager.Stop()
		err := os.Remove(dbFilePath)
		require.NoError(t, err)
	}()

	_, err := datastore.New(dbFilePath)
	assert.Error(t, err, "datastore manager should fail to open duplicate db file")
}

func TestManagerStop(t *testing.T) {
	manager, dbFilePath := initDatastoreManager(t)

	// verify that stopping it multiple times doesn't cause any issues
	manager.Stop()
	err := os.Remove(dbFilePath)
	require.NoError(t, err)
	manager.Stop()
}

func TestManagerGetPath(t *testing.T) {
	manager, dbFilePath := initDatastoreManager(t)
	defer func() {
		manager.Stop()
		err := os.Remove(dbFilePath)
		require.NoError(t, err)
	}()

	assert.Equal(t, dbFilePath, manager.GetPath())
}

func initDatastoreManager(t *testing.T) (*datastore.Manager, string) {
	t.Helper()

	ct := time.Now()
	fn := fmt.Sprintf("test_%d.db", ct.Unix())
	m, err := datastore.New(fn)
	require.NoError(t, err)

	appFS := afero.NewMemMapFs()

	// create test files and directories
	err = appFS.MkdirAll("test", 0o755)
	require.NoError(t, err)

	testValues := []testEntry{
		{
			Key:   1234,
			Value: "Mollit",
		},
		{
			Key:   1234,
			Value: "Comodo",
		},
		{
			Key:   5678,
			Value: "Conisere",
		},
	}

	for _, tv := range testValues {
		tv := tv // set value so we can pass to Save
		require.NoError(t, m.Save(&tv))
	}

	return m, fn
}

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestManager_Save(t *testing.T) {
	manager, dbFilePath := initDatastoreManager(t)

	defer func() {
		err := os.Remove(dbFilePath)
		require.NoError(t, err)
	}()

	type testData struct {
		ID      int
		TestNum int
	}

	td := testData{ID: 1, TestNum: 5}
	err := manager.Save(&td)
	assert.NoError(t, err)

	manager.Stop()
	td = testData{ID: 2, TestNum: 42}
	err = manager.Save(&td)
	assert.Error(t, err)
}

func TestManager_One(t *testing.T) {
	m, dbFilePath := initDatastoreManager(t)

	defer func() {
		m.Stop()
		err := os.Remove(dbFilePath)
		require.NoError(t, err)
	}()

	type args struct {
		fieldName string
		value     interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty search",
			args: args{
				fieldName: "",
				value:     nil,
			},
			wantErr: true,
		},
		{
			name: "good - with results",
			args: args{
				fieldName: "Key",
				value:     1234,
			},
			wantErr: false,
		},
		{
			name: "good - no results", // check for specific error "not found"
			args: args{
				fieldName: "Key",
				value:     5678,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got testEntry

			if err := m.One(tt.args.fieldName, tt.args.value, &got); (err != nil) != tt.wantErr {
				t.Errorf("Manager.One() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Find(t *testing.T) {
	m, dbFilePath := initDatastoreManager(t)

	defer func() {
		m.Stop()
		err := os.Remove(dbFilePath)
		require.NoError(t, err)
	}()

	type args struct {
		fieldName string
		value     interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "empty search",
			args: args{
				fieldName: "",
				value:     nil,
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "good - with results",
			args: args{
				fieldName: "Key",
				value:     1234,
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "good - no results", // check for specific error "not found"
			args: args{
				fieldName: "Key",
				value:     1,
			},
			wantCount: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []testEntry

			if err := m.Find(tt.args.fieldName, tt.args.value, &got); (err != nil) != tt.wantErr {
				t.Errorf("Manager.Find() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equalf(t, tt.wantCount, len(got), "Found %d results, wanted %d", len(got), tt.wantCount)
		})
	}
}

func TestManager_AllByIndex(t *testing.T) {
	m, dbFilePath := initDatastoreManager(t)

	defer func() {
		m.Stop()
		err := os.Remove(dbFilePath)
		require.NoError(t, err)
	}()

	type args struct {
		fieldName string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{

		{
			name: "good - with results",
			args: args{
				fieldName: "Key",
			},
			wantCount: 3,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []testEntry

			if err := m.AllByIndex(tt.args.fieldName, &got); (err != nil) != tt.wantErr {
				t.Errorf("Manager.AllByIndex() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equalf(t, tt.wantCount, len(got), "Found %d results, wanted %d", len(got), tt.wantCount)
		})
	}
}

func TestManager_All(t *testing.T) {
	m, dbFilePath := initDatastoreManager(t)

	defer func() {
		m.Stop()
		err := os.Remove(dbFilePath)
		require.NoError(t, err)
	}()

	tests := []struct {
		name      string
		wantCount int
		wantErr   bool
	}{
		{
			name: "good - with results",

			wantCount: 3,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []testEntry

			if err := m.All(&got); (err != nil) != tt.wantErr {
				t.Errorf("Manager.All() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equalf(t, tt.wantCount, len(got), "Found %d results, wanted %d", len(got), tt.wantCount)
		})
	}
}

func TestManager_Select(t *testing.T) {
	m, dbFilePath := initDatastoreManager(t)

	defer func() {
		m.Stop()
		err := os.Remove(dbFilePath)
		require.NoError(t, err)
	}()

	type args struct {
		matcher q.Matcher
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "good - with results",
			args: args{
				matcher: q.Eq("Key", 1234),
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "good - no results", // check for specific error "not found"
			args: args{
				matcher: q.Eq("Key", 1),
			},
			wantCount: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []testEntry

			q := m.Select(tt.args.matcher)
			if err := q.Find(&got); (err != nil) != tt.wantErr {
				t.Errorf("Manager.Select() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equalf(t, tt.wantCount, len(got), "Found %d results, wanted %d", len(got), tt.wantCount)
		})
	}
}

func TestManager_Range(t *testing.T) {
	m, dbFilePath := initDatastoreManager(t)

	defer func() {
		m.Stop()
		err := os.Remove(dbFilePath)
		require.NoError(t, err)
	}()

	type args struct {
		fieldName string
		min       interface{}
		max       interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "empty search",
			args: args{
				fieldName: "",
				min:       nil,
				max:       nil,
			},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "good - subset results",
			args: args{
				fieldName: "Key",
				min:       1235,
				max:       5679,
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name: "good - full results", // check for specific error "not found"
			args: args{
				fieldName: "Key",
				min:       1,
				max:       9999,
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name: "good - no results", // check for specific error "not found"
			args: args{
				fieldName: "Key",
				min:       nil,
				max:       nil,
			},
			wantCount: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []testEntry
			if err := m.Range(tt.args.fieldName, tt.args.min, tt.args.max, &got); (err != nil) != tt.wantErr {
				t.Errorf("Manager.Range() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equalf(t, tt.wantCount, len(got), "Found %d results, wanted %d", len(got), tt.wantCount)
		})
	}
}

func TestManager_Prefix(t *testing.T) {
	m, dbFilePath := initDatastoreManager(t)

	defer func() {
		m.Stop()
		err := os.Remove(dbFilePath)
		require.NoError(t, err)
	}()

	type args struct {
		fieldName string
		prefix    string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "good - with results",
			args: args{
				fieldName: "Value",
				prefix:    "Co",
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "good - no results", // check for specific error "not found"
			args: args{
				fieldName: "Value",
				prefix:    "Zzz",
			},
			wantCount: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []testEntry

			if err := m.Prefix(tt.args.fieldName, tt.args.prefix, &got); (err != nil) != tt.wantErr {
				t.Errorf("Manager.Prefix() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equalf(t, tt.wantCount, len(got), "Found %d results, wanted %d", len(got), tt.wantCount)
		})
	}
}

func TestManager_Count(t *testing.T) {
	m, dbFilePath := initDatastoreManager(t)

	defer func() {
		m.Stop()
		err := os.Remove(dbFilePath)
		assert.NoError(t, err)
	}()

	tests := []struct {
		name      string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "count all listings",
			wantCount: 3,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got testEntry

			c, err := m.Count(&got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.Find() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equalf(t, tt.wantCount, c, "Found %d results, wanted %d", c, tt.wantCount)
		})
	}
}
