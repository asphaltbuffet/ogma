// Package datastore is a simple package for tracking of a single instance of application data.
package datastore

import (
	"fmt"

	"github.com/asdine/storm/v3"
)

// Manager main object for data store.
type Manager struct {
	Store    *storm.DB
	filePath string
}

// New returns a new datastore Manager.
func New(filePath string) *Manager {
	storm, err := storm.Open(filePath)
	if err != nil {
		fmt.Println("Failed to open db: ", err)
	}

	return &Manager{
		Store:    storm,
		filePath: filePath,
	}
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
