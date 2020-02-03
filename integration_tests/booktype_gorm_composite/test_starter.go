package booktype_gorm_composite

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/merlinapp/datarepo-go"
	"github.com/merlinapp/datarepo-go/cachestore/composite"
	"github.com/merlinapp/datarepo-go/cachestore/memory"
	"github.com/merlinapp/datarepo-go/cachestore/nocache"
	"github.com/merlinapp/datarepo-go/cachestore/stats"
	"github.com/merlinapp/datarepo-go/integration_tests"
	"github.com/merlinapp/datarepo-go/integration_tests/booktype_gorm_composite/testdomain"
	"github.com/merlinapp/datarepo-go/integration_tests/model"
	"github.com/merlinapp/datarepo-go/repo/gorm"
	"time"
)

var testInstance *testdomain.SystemInstance

func startSystemForIntegrationTests() *testdomain.SystemInstance {
	if testInstance != nil {
		return testInstance
	}

	db := integration_tests.TestConnectionFactory()
	db.LogMode(true)

	cacheSize := 1 * 1024 * 1024
	memoryCacheStore := memory.NewFreeCacheInMemoryStore(cacheSize)
	emptyCacheStore := &nocache.Store{}

	bookTypeCacheStore := composite.NewCompositeCacheStore(emptyCacheStore, memoryCacheStore)
	bookTypeStatsCacheStore := stats.NewStatsCacheStore(bookTypeCacheStore)

	bookTypeRepo := gorm.CachedRepositoryBuilder(db, &model.BookType{}).
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
