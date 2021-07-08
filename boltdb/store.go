package boltdb

import "github.com/boltdb/bolt"

type Store struct {
	db      *bolt.DB
	FirstDB FirstDB
}

func Init(path string) *Store {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		panic(err)
	}

	return &Store{
		db:      db,
		FirstDB: newFirstDB(db),
	}
}

func (ins *Store) Close() error {
	return ins.db.Close()
}
