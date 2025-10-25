package storage

import "github.com/dgraph-io/badger/v3"

type KvScan struct {
	txn      *badger.Txn
	iterator *badger.Iterator
	prefix   []byte
	keyDst   []byte
	valueDst []byte
}

func NewKvScan(txn *badger.Txn, prefix []byte) *KvScan {
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = true
	opts.Prefix = prefix

	it := txn.NewIterator(opts)
	it.Seek(prefix)

	return &KvScan{
		txn:      txn,
		iterator: it,
	}
}

func (it *KvScan) Valid() bool {
	return it.iterator.ValidForPrefix(it.prefix)
}

func (it *KvScan) Next() {
	it.iterator.Next()
}

func (it *KvScan) Item() ([]byte, []byte, error) {
	item := it.iterator.Item()
	it.keyDst = item.KeyCopy(it.keyDst)
	v, err := item.ValueCopy(it.valueDst)
	if err != nil {
		return nil, nil, err
	}
	it.valueDst = v
	return it.keyDst, it.valueDst, nil
}

func (it *KvScan) Close() {
	it.iterator.Close()
	it.txn.Discard()
}
