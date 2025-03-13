package datasource

import (
	"errors"
	"sync"
)

type DataSource interface {
	Read(key string) (string, error)
	Create(key string, value interface{}) error
	Update(key string, value interface{}) error
	Upsert(key string, value interface{})
	Delete(key string) error
}

// InMemoryDB represents a thread-safe in-memory database using sync.Map.
type InMemoryDB struct {
	data sync.Map
}

// NewInMemoryDB initializes a new in-memory database.
func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		data: sync.Map{},
	}
}

// Create adds a new key-value pair to the database.
func (db *InMemoryDB) Create(key string, value interface{}) error {
	_, loaded := db.data.LoadOrStore(key, value)
	if loaded {
		return errors.New("key already exists")
	}
	return nil
}

// Read retrieves the value associated with the given key.
func (db *InMemoryDB) Read(key string) (string, error) {
	value, ok := db.data.Load(key)
	if !ok {
		return "", errors.New("key not found")
	}
	return value.(string), nil
}

// Upsert updates the value if the key exists, or inserts the key-value pair if it doesn't.
func (db *InMemoryDB) Upsert(key string, value interface{}) {
	db.data.Store(key, value) // Store will insert or update the key
}

// Update updates the value associated with the given key.
func (db *InMemoryDB) Update(key string, value interface{}) error {
	if _, ok := db.data.Load(key); !ok {
		return errors.New("key not found")
	}
	db.data.Store(key, value)
	return nil
}

// Delete removes the key-value pair associated with the given key.
func (db *InMemoryDB) Delete(key string) error {
	if _, ok := db.data.Load(key); !ok {
		return errors.New("key not found")
	}
	db.data.Delete(key)
	return nil
}
