package main

import (
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v2"
)

var ErrNotFound = errors.New("entry not found")

type db struct {
	store *badger.DB
}

func newDB(path string) (db, error) {
	badgerDB, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		fmt.Errorf("open failed: %w", err)
	}

	return db{store: badgerDB}, nil
}

func (db db) close() error {
	return db.close()
}

func (db db) add(source string, target string) error {
	return db.store.Update(func(tx *badger.Txn) error {
		err := tx.Set([]byte(source), []byte(target))
		if err != nil {
			return fmt.Errorf("set failed: %w", err)
		}
		return nil
	})
}

func (db db) delete(source string) error {
	return db.store.Update(func(tx *badger.Txn) error {
		err := tx.Delete([]byte(source))
		if err != nil {
			return fmt.Errorf("delete failed: %w", err)
		}
		return nil
	})
}

func (db db) get(source string) (target string, err error) {
	err = db.store.View(func(tx *badger.Txn) error {
		val, err := tx.Get([]byte(source))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return ErrNotFound
			}

			return err
		}
		valueCopy, err := val.ValueCopy(nil)
		if err != nil {
		    return fmt.Errorf("copy failed: %w", err)
		}

		target = string(valueCopy)
		return nil
	})

	return target, err
}

func (db db) list() (map[string]string, error) {
	redirects := map[string]string{}

	err := db.store.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				redirects[string(k)] = string(v);
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return redirects, err
}
