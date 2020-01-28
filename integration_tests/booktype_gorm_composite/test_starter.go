package booktype_gorm_composite

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/merlinapp/datarepo-go"
	"github.com/merlinapp/datarepo-go/cachestore/composite"
	"github.com/merlinapp/datarepo-go/cachestore/memory"
	"github.com/merlinapp/datarepo-go/cachestore/nocache"
	"github.com/merlinapp/datarepo-go/cachestore/stats"
	"github.com/merlinapp/datarepo-go/integration_tests"
	"github.com/merlinapp/datarepo-go/integration_tests/booktype_gorm_composite/testdomain"
	"github.com/merlinapp/datarepo-go/integration_tests/model"
	gorm2 "github.com/merlinapp/datarepo-go/repo/gorm"
	"time"
)

type SystemInstance struct {
	Ctx                context.Context
	DB                 *gorm.DB
	BookTypeCacheStore stats.StatsCacheStore
	BookTypeRepo       datarepo.CachedRepository
}

var testInstance *testdomain.SystemInstance

func startSystemForIntegrationTests() *testdomain.SystemInstance {
	if testInstance != nil {
		return testInstance
	}

	db := integration_tests.TestConnectionFactory()
	db.LogMode(true)

	memoryCacheStore := memory.NewFreeCacheInMemoryStore(1 * 1024 * 1024)
	emptyCacheStore := &nocache.Store{}

	bookTypeCacheStore := composite.NewCompositeCacheStore(emptyCacheStore, memoryCacheStore)
	bookTypeStatsCacheStore := stats.NewStatsCacheStore(bookTypeCacheStore)

	bookTypeRepo := gorm2.CachedRepositoryBuilder(db, &model.BookType{}).
		WithUniqueKeyCache(bookTypeCache, bookTypeStatsCacheStore).
		BuildCachedRepository()

	testInstance = &testdomain.SystemInstance{
		Ctx:                context.Background(),
		DB:                 db,
		BookTypeCacheStore: bookTypeStatsCacheStore,
		BookTypeRepo:       bookTypeRepo,
	}

	return testInstance
}

func prepareTestDB() {
	testInstance.DB.Delete(&model.BookType{})
}

func rollbackTestDb() {
	testInstance.DB.Close()
	testInstance = nil
}

var (
	bookTypeCache = datarepo.UniqueKeyCacheDefinition{
		KeyPrefix:    "bt:",
		KeyFieldName: "ID",
		Expiration:   5 * time.Minute,
	}
)
