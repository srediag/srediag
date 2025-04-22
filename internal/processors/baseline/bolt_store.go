package baseline

import (
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

var baselineBucket = []byte("baselines")

type boltStore struct {
	db *bbolt.DB
}

func newBoltStore(path string) (Store, error) {
	db, err := bbolt.Open(path, 0600, &bbolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open bolt db: %w", err)
	}

	// Create bucket if it doesn't exist
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(baselineBucket)
		return err
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create bucket: %w", err)
	}

	return &boltStore{db: db}, nil
}

func (s *boltStore) Get(path string) (string, error) {
	var hash string
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(baselineBucket)
		v := b.Get([]byte(path))
		if v == nil {
			return fmt.Errorf("hash not found for path: %s", path)
		}
		hash = string(v)
		return nil
	})
	return hash, err
}

func (s *boltStore) Set(path, hash string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(baselineBucket)
		return b.Put([]byte(path), []byte(hash))
	})
}

func (s *boltStore) Delete(path string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(baselineBucket)
		return b.Delete([]byte(path))
	})
}

func (s *boltStore) Close() error {
	return s.db.Close()
}
