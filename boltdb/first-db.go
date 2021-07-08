package boltdb

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

type FirstData struct {
	ID   uint64
	Data []byte
}

type firstdb struct {
	db         *bolt.DB
	bucketName []byte
	sync.Mutex
}

type FirstDB interface {
	Insert(item *FirstData) error
	Del(id uint64) error
	Get(id uint64) (FirstData, error)
	Load() error
}

func newFirstDB(db *bolt.DB) FirstDB {
	return &firstdb{
		db:         db,
		bucketName: bytes.NewBufferString("first-db").Bytes(),
	}
}

func (ins *firstdb) Insert(item *FirstData) error {
	ins.Lock()
	defer ins.Unlock()

	return ins.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(ins.bucketName)
		if err != nil {
			return err
		}
		if item.ID == 0 {
			item.ID, _ = bucket.NextSequence()
		} else if duplicate := bucket.Get(uitob(item.ID)); duplicate != nil {
			item.ID, _ = bucket.NextSequence()
		}
		var (
			gob_buffer bytes.Buffer
			key        = uitob(item.ID)
		)
		if err := gob.NewEncoder(&gob_buffer).Encode(item); err != nil {
			return err
		}
		return bucket.Put(key, gob_buffer.Bytes())
	})
}

func (ins *firstdb) Del(id uint64) error {
	ins.Lock()
	defer ins.Unlock()
	return ins.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(ins.bucketName)
		if bucket != nil {
			// pass
		} else {
			return errors.New("error no documents")
		}
		return bucket.Delete(uitob(id))
	})
}

func (ins *firstdb) Get(id uint64) (FirstData, error) {
	ins.Lock()
	defer ins.Unlock()
	var (
		item FirstData
	)
	if err := ins.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(ins.bucketName)
		if bucket != nil {
			// pass
		} else {
			return errors.New("error no documents")
		}
		v := bucket.Get(uitob(id))
		if len(v) > 0 {
			var (
				gob_buffer = bytes.NewReader(v)
			)
			if err := gob.NewDecoder(gob_buffer).Decode(&item); err != nil {
				return err
			}
			return nil
		}
		return errors.New("error no documents")
	}); err != nil {
		return FirstData{}, err
	}

	return item, nil
}

func (ins *firstdb) Load() error {
	ins.Lock()
	defer ins.Unlock()

	return ins.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(ins.bucketName)
		if bucket != nil {
			// pass
		} else {
			return errors.New("No documents")
		}
		return bucket.ForEach(func(k, v []byte) error {
			var (
				gob_buffer = bytes.NewReader(v)
				data       FirstData
			)
			if err := gob.NewDecoder(gob_buffer).Decode(&data); err != nil {
				return err
			}
			fmt.Printf("key=%d, value=%+v\n", obtui(k), data.Data)
			return nil
		})
	})
}
