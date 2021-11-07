package datastore_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
)

func TestManagerNew(t *testing.T) {
	manager, dbFilePath, err := initDatastoreManager()
	defer func() {
		err = os.Remove(dbFilePath)
		assert.NoError(t, err)
	}()
	defer manager.Stop()

	assert.NoError(t, err)

	_, err = os.Stat(dbFilePath)

	assert.NotErrorIs(t, err, os.ErrNotExist)
}

func TestManagerNewFail(t *testing.T) {
	manager, dbFilePath, err := initDatastoreManager()
	defer func() {
		err = os.Remove(dbFilePath)
		assert.NoError(t, err)
	}()
	defer manager.Stop()

	_, err = datastore.New(dbFilePath)
	assert.Error(t, err, "datastore manager should fail to open duplicate db file")
}

func TestManagerGetPath(t *testing.T) {
	manager, dbFilePath, err := initDatastoreManager()
	defer func() {
		err = os.Remove(dbFilePath)
		assert.NoError(t, err)
	}()
	defer manager.Stop()

	assert.Equal(t, dbFilePath, manager.GetPath())
}

func initDatastoreManager() (*datastore.Manager, string, error) {
	currentTime := time.Now()
	filename := fmt.Sprintf("test_%d.db", currentTime.Unix())
	manager, err := datastore.New(filename)

	return manager, filename, err
}

func init() {
	log.SetOutput(ioutil.Discard)
}
