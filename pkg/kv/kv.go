package kv

// This package abstracts the key-value database library used

import (
	"fmt"
	"log"

	"github.com/akrylysov/pogreb"
	"github.com/ebarped/kubeswap/pkg/kubeconfig"
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

func (kv *DB) Close() {
	// this should:
	// - compress the real database folder into a single file
	err := kv.db.Close()
	if err != nil {
		log.Fatalf("error closing db: %s", err)
	}
}

func (kv *DB) GetKubeconfig(key string) (*kubeconfig.Kubeconfig, error) {
	exist, err := kv.db.Has([]byte(key))
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("key does not exist in the db: %s", key)
	}

	val, err := kv.db.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	return &kubeconfig.Kubeconfig{
		Name:    key,
		Content: string(val),
	}, nil
}

func (kv *DB) PutKubeconfig(key string, value []byte) error {
	exist, err := kv.db.Has([]byte(key))
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("key already exists in the db: %s", key)
	}

	err = kv.db.Put([]byte(key), value)
	if err != nil {
		return err
	}
	return nil
}

func (kv *DB) DeleteKubeconfig(key string) error {
	exist, err := kv.db.Has([]byte(key))
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("key does not exist in the db: %s", key)
	}

	err = kv.db.Delete([]byte(key))
	if err != nil {
		return err
	}
	return nil
}

func (kv *DB) GetAll() ([]kubeconfig.Kubeconfig, error) {
	var items []kubeconfig.Kubeconfig
	it := kv.db.Items()
	for {
		key, val, err := it.Next()
		if err == pogreb.ErrIterationDone {
			break
		}
		if err != nil {
			return nil, err
		}
		kc := kubeconfig.Kubeconfig{
			Name:    string(key),
			Content: string(val),
		}
		items = append(items, kc)
	}
	return items, nil
}

func (kv *DB) IsEmpty() bool {
	return kv.db.Count() == 0
}
