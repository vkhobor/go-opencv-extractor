package db

import (
	"encoding/json"

	"github.com/go-errors/errors"
)

type DbInMemory struct {
	db map[string]json.RawMessage
}

func NewDbInMemory(rawDb string) (*DbInMemory, error) {
	db := make(map[string]json.RawMessage)
	err := json.Unmarshal([]byte(rawDb), &db)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return &DbInMemory{
		db: db,
	}, nil
}

func ReadAllEntries[T any](db *DbInMemory) (map[string]T, error) {
	entries := make(map[string]T)
	for key, value := range db.db {
		var e T
		err := json.Unmarshal(value, &e)
		if err != nil {
			return nil, errors.Wrap(err, 0)
		}
		entries[key] = e
	}
	return entries, nil
}

func ReadEntry[T any](db *DbInMemory, key string) (bool, T, error) {
	var e T
	value, ok := db.db[key]
	if !ok {
		return false, e, nil
	}

	err := json.Unmarshal(value, &e)
	if err != nil {
		return false, e, errors.Wrap(err, 0)
	}
	return true, e, nil
}

func WriteEntry[T any](db *DbInMemory, key string, value T) error {
	raw, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, 0)
	}
	db.db[key] = raw
	return nil
}

func (db *DbInMemory) String() ([]byte, error) {
	return json.MarshalIndent(db.db, "", "   ")
}
