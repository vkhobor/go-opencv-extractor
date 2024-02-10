package db

import (
	"context"
	"os"
	"time"

	"github.com/go-errors/errors"

	"github.com/gofrs/flock"
)

type Db[T any] struct {
	dbPath     string
	dbLockPath string
	dbLock     *flock.Flock
}

func OpenDb[T any](dbPath string, lockPath string) (*Db[T], error) {
	db := Db[T]{
		dbPath:     dbPath,
		dbLockPath: lockPath,
		dbLock:     flock.New(lockPath),
	}

	// Create timeout context
	ctx, _ := context.WithTimeout(context.Background(), time.Second*45)
	err := db.ensureExists(ctx)

	return &db, err
}

func (db *Db[T]) ensureExists(ctx context.Context) error {
	_, err := db.dbLock.TryRLockContext(ctx, time.Second*3)
	if err != nil {
		return errors.Wrap(err, 0)
	}
	defer db.dbLock.Close()

	// Check if the file exists
	_, err = os.Stat(db.dbPath)
	if os.IsNotExist(err) {
		// Create the file if it doesn't exist
		file, err := os.Create(db.dbPath)
		if err != nil {
			return errors.Wrap(err, 0)
		}
		defer file.Close()

		// Write the file structure as basic json object
		_, err = file.WriteString("{}")
		if err != nil {
			return errors.Wrap(err, 0)
		}
	}
	return nil
}

// Put entry into db fail if lock is held
func (db *Db[T]) TryPut(key string, value T) error {
	ok, err := db.dbLock.TryLock()
	if err != nil {
		return errors.Wrap(err, 0)
	}
	if !ok {
		return errors.Errorf("lock is held")
	}

	defer db.dbLock.Close()

	return db.put(key, value)
}

// Put entry into db block until lock is released
func (db *Db[T]) Put(ctx context.Context, key string, value T) error {
	_, err := db.dbLock.TryLockContext(ctx, time.Second*3)
	if err != nil {
		return errors.Wrap(err, 0)
	}
	defer db.dbLock.Close()

	return db.put(key, value)
}

// unguarded put
func (db *Db[T]) put(key string, value T) error {
	file, err := os.ReadFile(db.dbPath)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	fileString := string(file)
	dbMemory, err := NewDbInMemory(fileString)
	if err != nil {
		return err
	}

	err = WriteEntry(dbMemory, key, value)
	if err != nil {
		return err
	}

	newFile, err := dbMemory.String()
	if err != nil {
		return err
	}

	err = os.WriteFile(db.dbPath, newFile, 0644)
	if err != nil {
		return errors.Wrap(err, 0)
	}

	return nil
}

func (db *Db[T]) ReadAll(ctx context.Context) (map[string]T, error) {
	_, err := db.dbLock.TryRLockContext(ctx, time.Second*3)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}
	defer db.dbLock.Close()

	file, err := os.ReadFile(db.dbPath)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	fileString := string(file)
	dbMemory, err := NewDbInMemory(fileString)
	if err != nil {
		return nil, errors.Wrap(err, 0)
	}

	return ReadAllEntries[T](dbMemory)
}

// Read entry from db
func (db *Db[T]) Read(ctx context.Context, key string) (bool, T, error) {
	_, err := db.dbLock.TryRLockContext(ctx, time.Second*3)
	if err != nil {
		return false, getZero[T](), errors.Wrap(err, 0)
	}
	defer db.dbLock.Close()

	file, err := os.ReadFile(db.dbPath)
	if err != nil {
		return false, getZero[T](), errors.Wrap(err, 0)
	}

	fileString := string(file)
	dbMemory, err := NewDbInMemory(fileString)
	if err != nil {
		return false, getZero[T](), errors.Wrap(err, 0)
	}

	return ReadEntry[T](dbMemory, key)
}

func getZero[T any]() T {
	var result T
	return result
}
