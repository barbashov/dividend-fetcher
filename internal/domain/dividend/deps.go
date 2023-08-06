package dividend

type Cache interface {
	Set(key string, value []Dividends) error
	Get(key string) ([]Dividends, error)
	Close() error
}
