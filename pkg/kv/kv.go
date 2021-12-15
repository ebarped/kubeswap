package kv

// This package abstracts the key-value database library used

import (
	"github.com/akrylysov/pogreb"
	"github.com/ebarped/kubeswap/pkg/kubeconfig"
)

type DB struct {
	*pogreb.DB
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

func (kv *DB) PutKubeconfig(key string, value []byte) error {
	err := kv.Put([]byte(key), value)
	if err != nil {
		return err
	}
	return nil
}

func (kv *DB) CloseDB() {
	// this should:
	// - compress the real database folder into a single file
	kv.Close()
}

func (kv *DB) GetKubeconfig(key string) (*kubeconfig.Kubeconfig, error) {
	val, err := kv.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	return &kubeconfig.Kubeconfig{
		Name:    key,
		Content: string(val),
	}, nil
}

func (kv *DB) GetAll() ([]kubeconfig.Kubeconfig, error) {
	var items []kubeconfig.Kubeconfig
	it := kv.Items()
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
