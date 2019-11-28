package stats

import (
	"context"
	"github.com/merlinapp/datarepo-go"
	"time"
)

type StatsCacheStore interface {
	datarepo.CacheStore
	ClearStats()
	// number of Set operations performed in the CacheStore
	Sets() int64
	// number of Cache Hits
	Hits() int64
	// number of Cache Misses
	Miss() int64
	// number of Delete operations performed in the CacheStore
	Dels() int64
}

type statsCacheStore struct {
	delegate datarepo.CacheStore
	sets     int64
	hits     int64
	miss     int64
	dels     int64
}

// Creates a new CacheStore that keeps hit/miss/set/del stats for an underlying/delegate
// CacheStore
func NewStatsCacheStore(delegate datarepo.CacheStore) StatsCacheStore {
	store := statsCacheStore{delegate: delegate}
	return &store
}

func (s *statsCacheStore) Delete(ctx context.Context, key string) error {
	err := s.delegate.Delete(ctx, key)
	if err != nil {
		return err
	}
	s.dels++
	return nil
}

func (s *statsCacheStore) Get(ctx context.Context, key string, out interface{}) (bool, error) {
	found, err := s.delegate.Get(ctx, key, out)
	if err != nil {
		return found, err
	}

	if found {
		s.hits++
	} else {
		s.miss++
	}
	return found, err
}

func (s *statsCacheStore) GetMulti(ctx context.Context, keys []string, out interface{}) ([]bool, error) {
	foundArr, err := s.delegate.GetMulti(ctx, keys, out)
	if err != nil {
		return foundArr, err
	}

	for _, v := range foundArr {
		if v {
			s.hits++
		} else {
			s.miss++
		}
	}
	return foundArr, err
}

func (s *statsCacheStore) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	s.delegate.Set(ctx, key, value, expiration)
	s.sets++
}

func (s *statsCacheStore) ClearStats() {
	s.hits = 0
	s.miss = 0
	s.dels = 0
	s.sets = 0
}

func (s *statsCacheStore) Sets() int64 {
	return s.sets
}

func (s *statsCacheStore) Hits() int64 {
	return s.hits
}

func (s *statsCacheStore) Miss() int64 {
	return s.miss
}

func (s *statsCacheStore) Dels() int64 {
	return s.dels
}
