// Package datastore is a simple package for tracking of a single instance of application data.
package datastore

import (
	"fmt"

	"github.com/asdine/storm/v3"
	log "github.com/sirupsen/logrus"
)

// Manager main object for data store.
type Manager struct {
	Store    *storm.DB
	filePath string
}

// New returns a new datastore Manager.
func New(filePath string) (*Manager, error) {
	// TODO: can this be converted to use afero? 2021-11-07 BL
	storm, err := storm.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open db: ", err)
		log.WithFields(log.Fields{
			"filePath": filePath,
		}).Error("Failed to open db.")

		return nil, err
	}

	return &Manager{
		Store:    storm,
		filePath: filePath,
	}, nil
}

// GetPath returns the filepath to db file.
func (m *Manager) GetPath() string {
	return m.filePath
}

// Stop stops database and any associated goroutines.
func (m *Manager) Stop() {
	if err := m.Store.Close(); err != nil {
		fmt.Println("Failed to close store: ", err)
	}
}

// Save saves data into the datastore.
func (m *Manager) Save(data interface{}) error {
	if err := m.Store.Save(data); err != nil {
		return err
	}

	return nil
}

// Data provides access to datastore storage.
func (m *Manager) Data() *storm.DB {
	return m.Store
}

// A Writer can write to a datastore.
type Writer interface {
	Save(data interface{}) error
}

// A Reader can read from a datastore.
type Reader interface {
	Data() *storm.DB
}
