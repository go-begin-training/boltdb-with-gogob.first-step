package main

import (
	"bytes"
	"fmt"
	"math"

	"github.com/go-begin-training/boltdb-with-gogob.first-step/boltdb"
)

var (
	store_path string = "./store.bwg"
)

func main() {
	var (
		boltdb_service = boltdb.Init(store_path)
	)

	defer boltdb_service.Close()

	for id := uint64(0); id <= 20777; id++ {
		item := boltdb.FirstData{
			ID:   id,
			Data: bytes.NewBufferString("compressed").Bytes(),
		}
		if err := boltdb_service.FirstDB.Insert(&item); err != nil {
			fmt.Println(err)
		}

		fmt.Println("inserted id=", item.ID)
	}

	for id := uint64(0); id <= 20777; id++ {
		if math.Mod(float64(id), 2) == 0 {
			fmt.Printf("key=%d, deleting...\n", id)
			if err := boltdb_service.FirstDB.Del(id); err != nil {
				fmt.Printf("key=%d, error: %+v\n", id, err)
			}
		}
	}

	for id := uint64(0); id <= 20777; id++ {
		data, err := boltdb_service.FirstDB.Get(id)
		if err != nil {
			fmt.Printf("key=%d, error: %+v\n", id, err.Error())
			continue
		}
		fmt.Printf("key=%d, value=%+v\n", id, data.Data)

	}

}
