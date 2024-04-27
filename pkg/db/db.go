package db

import (
	"log"
	"time"

	"github.com/tidwall/buntdb"
)

type KV struct {
	db *buntdb.DB
}

func NewBuntDb(pathToDb string) (*KV, error) {
	db, err := buntdb.Open(pathToDb)
	if err != nil {
		log.Fatal(err)
	}
	return &KV{db: db}, nil
}

func (k *KV) Close() error {
	return k.db.Close()
}

func (k *KV) CreateJsonIndex(index string) {
	k.db.CreateIndex(index, "*", buntdb.IndexJSON(index))
}

func (k *KV) Add(key string, value []byte) error {
	err := k.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, string(value), nil)
		return err
	})
	if err != nil {
		panic(err)
	}
	return nil
}

func (k *KV) Get(key string) (string, error) {
	var obj string
	err := k.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(key)
		if err != nil {
			return err
		}
		obj = val
		return nil
	})
	return obj, err
}

func (k *KV) GetDescendRange(index, lessOrEqual, greaterThan string) ([]string, error) {
	var objs []string
	err := k.db.View(func(tx *buntdb.Tx) error {
		return tx.DescendRange(index, lessOrEqual, greaterThan, func(key, value string) bool {
			objs = append(objs, value)
			return true
		})
	})
	return objs, err
}

func (k *KV) Delete(key string) error {
	// buntdb do not support delete command, so we use set with ttl
	err := k.db.Update(func(tx *buntdb.Tx) error {
		tx.Set(key, "", &buntdb.SetOptions{Expires: true, TTL: time.Second})
		return nil
	})
	return err
}
