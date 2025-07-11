package kv

// This package abstracts the key-value database library used

import (
	"fmt"
	"log"

	"github.com/ebarped/kubeswap/pkg/kubeconfig"
	"go.etcd.io/bbolt"
)

const mainBucket = "Store"

type DB struct {
	db *bbolt.DB
}

func Open(path string) (*DB, error) {
	db, err := bbolt.Open(path, 0o600, nil)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (kv *DB) Close() {
	err := kv.db.Close()
	if err != nil {
		log.Fatalf("error closing db: %s", err)
	}
}

func (kv *DB) GetKubeconfig(key string) (*kubeconfig.Kubeconfig, error) {
	var val []byte

	err := kv.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(mainBucket))
		if b == nil {
			return fmt.Errorf("bucket not found, you have to sync first")
		}
		val = b.Get([]byte(key))
		if val == nil {
			return fmt.Errorf("key does not exists in the db: %s", key)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &kubeconfig.Kubeconfig{
		Name:    key,
		Content: string(val),
	}, nil
}

func (kv *DB) PutKubeconfig(key string, value []byte) error {
	err := kv.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(mainBucket))
		if err != nil {
			return err
		}

		val := b.Get([]byte(key))
		if val != nil {
			return fmt.Errorf("key already exists in the db: %s", key)
		}

		err = b.Put([]byte(key), value)
		if err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (kv *DB) DeleteKubeconfig(key string) error {
	err := kv.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(mainBucket))
		if err != nil {
			return err
		}

		val := b.Get([]byte(key))
		if val == nil {
			return fmt.Errorf("key does not exists: %s", key)
		}

		err = b.Delete([]byte(key))
		if err != nil {
			return fmt.Errorf("error trying to delete key: %s. %s", key, err)
		}

		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func (kv *DB) GetAll() ([]kubeconfig.Kubeconfig, error) {
	var items []kubeconfig.Kubeconfig

	err := kv.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(mainBucket))
		if b == nil {
			return fmt.Errorf("bucket not found, you have to sync first")
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			kc := kubeconfig.Kubeconfig{
				Name:    string(k),
				Content: string(v),
			}
			items = append(items, kc)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (kv *DB) IsEmpty() bool {
	var empty bool

	err := kv.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(mainBucket))
		if b == nil {
			return fmt.Errorf("bucket not found, you have to sync first")
		}

		c := b.Cursor()
		k, _ := c.First()
		empty = (k == nil)
		return nil
	})
	if err != nil {
		log.Fatalf("error checking if DB is empty: %s\n", err)
	}

	return empty
}
