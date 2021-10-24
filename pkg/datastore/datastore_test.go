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
	manager, dbFilePath := initDatastoreManager()
	defer os.Remove(dbFilePath)
	defer manager.Stop()

	assert.NotNil(t, manager)

	_, err := os.Stat(dbFilePath)

	assert.False(t, os.IsNotExist(err))
}

func TestManagerStop(t *testing.T) {
	manager, dbFilePath := initDatastoreManager()
	defer os.Remove(dbFilePath)

	assert.NotNil(t, manager)

	manager.Stop()

	assert.Nil(t, manager.Store)
}

func initDatastoreManager() (*datastore.Manager, string) {
	currentTime := time.Now()
	filename := fmt.Sprintf("test_%d.db", currentTime.Unix())
	manager := datastore.New(filename)

	return manager, filename
}
