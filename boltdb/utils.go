package boltdb

import "encoding/binary"

func uitob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func obtui(v []byte) uint64 {
	return binary.BigEndian.Uint64(v)
}
