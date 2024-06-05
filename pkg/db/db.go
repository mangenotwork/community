package db

import (
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

type LocalDB struct {
	path  string
	table string
	Conn  *bolt.DB
}

func (ldb *LocalDB) Open() *bolt.DB {
	if _, err := os.Stat(ldb.path); os.IsNotExist(err) {
		file, _ := os.Create(ldb.path)
		_ = file.Close()
	}
	ldb.Conn, _ = bolt.Open(ldb.path, 0600, nil)
	return ldb.Conn
}

func (ldb *LocalDB) Set(key string, value []byte) error {
	ldb.Open()
	defer func() {
		_ = ldb.Conn.Close()
	}()
	return ldb.Conn.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ldb.table))
		if b == nil {
			return fmt.Errorf("table is null")
		}
		return b.Put([]byte(key), value)
	})
}

var InformationTable = &LocalDB{
	path:  "./information",
	table: "information",
}
