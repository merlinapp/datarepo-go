package stats

import (
	"context"
	"github.com/merlinapp/datarepo-go"
)

type StatsDataFetcher interface {
	datarepo.DataFetcher
	ClearStats()
	// number of ids sent for reading to the DataFetcher
	Reads() int64
}

type statsDataFetcher struct {
	delegate datarepo.DataFetcher
	reads    int64
}

// Creates a new DataFetcher that keeps stats for an underlying/delegate
// DataFetcher
func NewStatsDataFetcher(delegate datarepo.DataFetcher) StatsDataFetcher {
	store := statsDataFetcher{delegate: delegate}
	return &store
}

func (s *statsDataFetcher) FindByKey(ctx context.Context, keyFieldName string, id interface{}) (datarepo.Result, error) {
	s.reads++
	return s.delegate.FindByKey(ctx, keyFieldName, id)
}

func (s *statsDataFetcher) FindByKeys(ctx context.Context, keyFieldName string, ids []interface{}) ([]datarepo.Result, error) {
	s.reads += int64(len(ids))
	return s.delegate.FindByKeys(ctx, keyFieldName, ids)
}

func (s *statsDataFetcher) ClearStats() {
	s.reads = 0
}

func (s *statsDataFetcher) Reads() int64 {
	return s.reads
}
