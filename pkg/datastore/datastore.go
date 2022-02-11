// Package datastore is a simple package for tracking of a single instance of application data.
package datastore

import (
	"fmt"
	"os"

	storm "github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/index"
	"github.com/asdine/storm/v3/q"
	log "github.com/sirupsen/logrus"
)

//go:generate mockery --output=../../mocks --log-level=warn --name=Saver
// A Saver can write to a datastore.
type Saver interface {
	Save(data interface{}) error
	Begin(writable bool) (storm.Node, error)
}

//go:generate mockery --output=../../mocks --log-level=warn --name=SaveStopper
// A SaveStopper can write to or close a datastore.
type SaveStopper interface {
	Save(data interface{}) error
	Begin(writable bool) (storm.Node, error)
	Stop()
}

// A Finder can fetch types from BoltDB.
type Finder interface {
	storm.Finder
}

// Manager main object for data store.
type Manager struct {
	Store    *storm.DB
	filePath string
}

// New returns a new datastore Manager.
func New(filePath string) (*Manager, error) {
	storm, err := storm.Open(filePath)
	if err != nil {
		log.WithFields(log.Fields{
			"filePath": filePath,
		}).Error("error opening datastore file")

		return nil, fmt.Errorf("error opening datastore file: %w", err)
	}

	return &Manager{
		Store:    storm,
		filePath: filePath,
	}, nil
}

// Open returns a datastore Manager. Error if datastore file does not exist.
func Open(fp string) (*Manager, error) {
	if _, err := os.Stat(fp); err != nil {
		log.WithFields(log.Fields{
			"filePath": fp,
		}).Error("error accessing datastore file: ", err)
		return nil, fmt.Errorf("error accessing datastore file: %w", err)
	}

	return New(fp)
}

// Begin starts a transactional datastore instance.
func (m *Manager) Begin(writable bool) (storm.Node, error) {
	return m.Store.Begin(writable)
}

// GetPath returns the filepath to db file.
func (m *Manager) GetPath() string {
	return m.filePath
}

// Stop stops database and any associated goroutines.
func (m *Manager) Stop() {
	if err := m.Store.Close(); err != nil {
		log.Error("Failed to close store: ", err)
	}
}

// Save saves data into the datastore.
func (m *Manager) Save(data interface{}) error {
	if err := m.Store.Save(data); err != nil {
		log.WithFields(log.Fields{
			"record": data,
		}).Error("error saving record: ", err)
		return fmt.Errorf("error saving record=%+v: %w", data, err)
	}

	return nil
}

// One returns one record by the specified index.
func (m *Manager) One(fieldName string, value interface{}, to interface{}) error {
	return m.Store.One(fieldName, value, to)
}

// Find returns one or more records by the specified index.
func (m *Manager) Find(fieldName string, value interface{}, to interface{}, options ...func(*index.Options)) error {
	return m.Store.Find(fieldName, value, to, options...)
}

// AllByIndex gets all the records of a bucket that are indexed in the specified index.
func (m *Manager) AllByIndex(fieldName string, to interface{}, options ...func(*index.Options)) error {
	return m.Store.AllByIndex(fieldName, to, options...)
}

// All gets all the records of a bucket.
// If there are no records it returns no error and the 'to' parameter is set to an empty slice.
func (m *Manager) All(to interface{}, options ...func(*index.Options)) error {
	return m.Store.All(to, options...)
}

// Select a list of records that match a list of matchers. Doesn't use indexes.
func (m *Manager) Select(matchers ...q.Matcher) storm.Query {
	return m.Store.Select(matchers...)
}

// Range returns one or more records by the specified index within the specified range.
func (m *Manager) Range(fieldName string, min interface{}, max interface{}, to interface{}, options ...func(*index.Options)) error {
	return m.Store.Range(fieldName, min, max, to, options...)
}

// Prefix returns one or more records whose given field starts with the specified prefix.
func (m *Manager) Prefix(fieldName string, prefix string, to interface{}, options ...func(*index.Options)) error {
	return m.Store.Prefix(fieldName, prefix, to, options...)
}

// Count counts all the records of a bucket.
func (m *Manager) Count(data interface{}) (int, error) {
	return m.Store.Count(data)
}
