package stats

import (
	"context"
	"github.com/merlinapp/datarepo-go"
)

type StatsDataWriter interface {
	datarepo.DataWriter
	ClearStats()
	// number of Create operations invoked in the DataWriter
	Creates() int64
	// number of Update operations invoked in the DataWriter
	Updates() int64
}

type statsDataWriter struct {
	delegate datarepo.DataWriter
	creates  int64
	updates  int64
}

// Creates a new DataWriter that keeps stats for an underlying/delegate
// DataWriter
func NewStatsDataWriter(delegate datarepo.DataWriter) StatsDataWriter {
	store := statsDataWriter{delegate: delegate}
	return &store
}

func (s *statsDataWriter) Create(ctx context.Context, value interface{}) error {
	s.creates++
	return s.delegate.Create(ctx, value)
}

func (s *statsDataWriter) Update(ctx context.Context, value interface{}) error {
	s.updates++
	return s.delegate.Update(ctx, value)
}

func (s *statsDataWriter) ClearStats() {
	s.creates = 0
	s.updates = 0
}

func (s *statsDataWriter) Creates() int64 {
	return s.creates
}

func (s *statsDataWriter) Updates() int64 {
	return s.updates
}
