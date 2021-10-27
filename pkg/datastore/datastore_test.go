package datastore_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/asphaltbuffet/ogma/pkg/datastore"
)

func TestManagerNew(t *testing.T) {
	manager, dbFilePath, err := initDatastoreManager()
	defer func() {
		err = os.Remove(dbFilePath)
		assert.Nil(t, err)
	}()
	defer manager.Stop()

	assert.Nil(t, err)

	_, err = os.Stat(dbFilePath)

	assert.False(t, os.IsNotExist(err))
}

func TestManagerNewFail(t *testing.T) {
	manager, dbFilePath, err := initDatastoreManager()
	defer func() {
		err = os.Remove(dbFilePath)
		assert.Nil(t, err)
	}()
	defer manager.Stop()

	_, err = datastore.New(dbFilePath)
	assert.NotNil(t, err)
}

func initDatastoreManager() (*datastore.Manager, string, error) {
	currentTime := time.Now()
	filename := fmt.Sprintf("test_%d.db", currentTime.Unix())
	manager, err := datastore.New(filename)

	return manager, filename, err
}
