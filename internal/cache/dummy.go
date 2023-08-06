package cache

import "smartlab-dividend-fetcher/internal/domain/dividend"

type DummyCache struct {
}

func (c *DummyCache) Set(key string, value []dividend.Dividends) error {
	return nil
}

func (c *DummyCache) Get(key string) ([]dividend.Dividends, error) {
	return nil, ErrNoCache
}

func (c *DummyCache) Close() error {
	return nil
}
