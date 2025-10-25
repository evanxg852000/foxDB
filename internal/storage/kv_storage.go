package storage

import (
	"os"

	"github.com/dgraph-io/badger/v3"
)

type KvStorage struct {
	db *badger.DB
}

func NewKvStorage(dbPath string) (*KvStorage, error) {
	opts := badger.DefaultOptions(dbPath).WithLogger(nil) // Disable Badger logging
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	return &KvStorage{db: db}, nil
}

func (s *KvStorage) Close() error {
	return s.db.Close()
}

func (s *KvStorage) Sync() error {
	return s.db.Sync()
}

func (s *KvStorage) Remove() error {
	_ = s.db.Close()
	return os.RemoveAll(s.db.Opts().Dir)
}

func (s *KvStorage) Set(key, value []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (s *KvStorage) Get(key []byte) ([]byte, error) {
	var valCopy []byte
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		valCopy = val
		return nil
	})
	if err != nil {
		return nil, err
	}
	return valCopy, nil
}

func (s *KvStorage) Delete(key []byte) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func (s *KvStorage) Batch(fn func(txn *badger.Txn) error) error {
	return s.db.Update(fn)
}

func (s *KvStorage) Scan(prefix []byte) *KvScan {
	return NewKvScan(s.db.NewTransaction(false), prefix)
}
