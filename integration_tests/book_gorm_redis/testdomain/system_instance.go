package testdomain

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/merlinapp/datarepo-go"
	"github.com/merlinapp/datarepo-go/cachestore/stats"
	stats2 "github.com/merlinapp/datarepo-go/repo/stats"
)

type SystemInstance struct {
	Ctx                     context.Context
	DB                      *gorm.DB
	BookCacheStore          stats.StatsCacheStore
	UniqueKeyDataFetcher    stats2.StatsDataFetcher
	NonUniqueKeyDataFetcher stats2.StatsDataFetcher
	BookRepo                datarepo.CachedRepository
}