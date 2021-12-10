package carbon

import (
	"path"
	"time"

	"github.com/dgraph-io/badger/v3"
)

// Cache
type Cache struct {
	db *badger.DB
}

// NewCache
func NewCache(dir string) (*Cache, error) {
	opt := badger.DefaultOptions(dir)
	opt.ValueDir = path.Join(dir, "data")
	opt.Logger = nil

	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	// Badger GC
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
		again:
			err := db.RunValueLogGC(0.7)
			if err == nil {
				goto again
			}
		}
	}()

	return &Cache{
		db: db,
	}, nil
}

// Set a key,val with a ttl. if ttl is 0 or less, ttl then its disabled.
func (c *Cache) Set(key string, val []byte, ttl time.Duration) error {
	return c.db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), val)

		if ttl > 0 {
			e.WithTTL(ttl)
		}

		return txn.SetEntry(e)
	})
}

// Get
func (c *Cache) Get(key string) ([]byte, error) {
	var val []byte
	if err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(v []byte) error {
			val = v
			return nil
		})
	}); err != nil {
		return nil, err
	}

	return val, nil
}

// Del
func (c *Cache) Del(key string) error {
	return c.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

// ForEach - iterate over all keys if prefix == "", or a specific prefix
func (c *Cache) ForEach(prefix string, fn func(key string, val []byte) error) error {
	switch prefix {
	case "":
		return c.db.Update(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			for it.Rewind(); it.Valid(); it.Next() {
				item := it.Item()
				if err := item.Value(func(v []byte) error {
					return fn(string(item.Key()), v)
				}); err != nil {
					return err
				}
			}

			return nil
		})
	default:
		return c.db.Update(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()
			seek := []byte(prefix)
			for it.Seek(seek); it.ValidForPrefix(seek); it.Next() {
				item := it.Item()

				if err := item.Value(func(v []byte) error {
					return fn(string(item.Key()), v)
				}); err != nil {
					return err
				}
			}

			return nil
		})
	}
}

// Ping
func (c *Cache) Ping() error {
	return c.db.View(func(txn *badger.Txn) error {
		return txn.Commit()
	})
}

// Close
func (c *Cache) Close() {
	c.db.Close()
}
