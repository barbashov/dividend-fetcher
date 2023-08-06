package cache

import (
	"encoding/json"
	"errors"
	"smartlab-dividend-fetcher/internal/domain/dividend"
	"time"

	"github.com/boltdb/bolt"
)

const (
	DefaultBucket = "div-fetcher-cache"
	DefaultTTL    = time.Hour * 24 * 7 // 1 week
)

var ErrNoCache = errors.New("no cache")

type Cache struct {
	bucketName []byte
	ttl        time.Duration

	db *bolt.DB
}

type Option func(c *Cache)

func BucketName(n string) Option {
	return func(c *Cache) {
		c.bucketName = []byte(n)
	}
}

func TTL(ttl time.Duration) Option {
	return func(c *Cache) {
		c.ttl = ttl
	}
}

type cacheItem struct {
	Created time.Time
	Value   []dividend.Dividends
}

func NewCache(filename string, opts ...Option) (*Cache, error) {
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		return nil, err
	}

	ret := &Cache{
		bucketName: []byte(DefaultBucket),
		ttl:        DefaultTTL,
		db:         db,
	}

	for _, opt := range opts {
		opt(ret)
	}

	return ret, nil
}

func (c *Cache) Close() error {
	return c.db.Close()
}

func (c *Cache) Set(key string, value []dividend.Dividends) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		item := cacheItem{
			Created: time.Now(),
			Value:   value,
		}
		data, err := json.Marshal(item)
		if err != nil {
			return err
		}

		b, err := tx.CreateBucketIfNotExists(c.bucketName)
		if err != nil {
			return err
		}

		return b.Put([]byte(key), data)
	})
}

func (c *Cache) Get(key string) ([]dividend.Dividends, error) {
	item := cacheItem{}
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.bucketName)
		if b == nil {
			return ErrNoCache
		}

		data := b.Get([]byte(key))
		if len(data) == 0 {
			return ErrNoCache
		}

		if err := json.Unmarshal(data, &item); err != nil {
			return err
		}

		if time.Since(item.Created) > c.ttl {
			return ErrNoCache
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return item.Value, nil
}
