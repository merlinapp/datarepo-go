package testdomain

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/merlinapp/datarepo-go"
	"github.com/merlinapp/datarepo-go/cachestore/stats"
)

type SystemInstance struct {
	Ctx                    context.Context
	DB                     *gorm.DB
	BookCategoryCacheStore stats.StatsCacheStore
	BookCategoryRepo       datarepo.CachedRepository
}
