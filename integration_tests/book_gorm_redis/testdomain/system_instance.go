package testdomain

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/merlinapp/datarepo-go"
	"github.com/merlinapp/datarepo-go/cachestore/stats"
	statsrepo "github.com/merlinapp/datarepo-go/repo/stats"
)

type SystemInstance struct {
	Ctx                     context.Context
	DB                      *gorm.DB
	BookCacheStore          stats.StatsCacheStore
	UniqueKeyDataFetcher    statsrepo.StatsDataFetcher
	NonUniqueKeyDataFetcher statsrepo.StatsDataFetcher
	BookRepo                datarepo.CachedRepository
}
