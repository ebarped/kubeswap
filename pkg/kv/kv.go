package kv

// This package abstracts the key-value database library used

import (
	"github.com/akrylysov/pogreb"
)

type DB struct {
	db *pogreb.DB
}

func Open(path string) (*DB, error) {
	// this should:
	// - take a compressed file and decompress it in a folder in temporal location
	// - use that location as the real db
	db, err := pogreb.Open(path, nil)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (kv *DB) Put(key string, value []byte) error {
	err := kv.db.Put([]byte(key), value)
	if err != nil {
		return err
	}
	return nil
}

func (kv *DB) Get(key string) (string, error) {
	val, err := kv.db.Get([]byte(key))
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func (kv *DB) CloseDB() {
	// this should:
	// - compress the real database folder into a single file
	kv.db.Close()
}
