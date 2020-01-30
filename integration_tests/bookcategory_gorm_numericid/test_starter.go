package bookcategory_gorm_memory

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/merlinapp/datarepo-go"
	"github.com/merlinapp/datarepo-go/cachestore/memory"
	"github.com/merlinapp/datarepo-go/cachestore/stats"
	"github.com/merlinapp/datarepo-go/integration_tests"
	"github.com/merlinapp/datarepo-go/integration_tests/bookcategory_gorm_numericid/testdomain"
	"github.com/merlinapp/datarepo-go/integration_tests/model"
	gorm2 "github.com/merlinapp/datarepo-go/repo/gorm"
	"time"
)

var testInstance *testdomain.SystemInstance

func startSystemForIntegrationTests() *testdomain.SystemInstance {
	if testInstance != nil {
		return testInstance
	}

	db := integration_tests.TestConnectionFactory()
	db.LogMode(true)

	cacheStore := memory.NewFreeCacheInMemoryStore(1 * 1024 * 1024)
	statsCacheStore := stats.NewStatsCacheStore(cacheStore)

	bookCategoryRepo := gorm2.CachedRepositoryBuilder(db, &model.BookCategory{}).
		WithUniqueKeyCache(bookCategoryCache, statsCacheStore).
		BuildCachedRepository()

	testInstance = &testdomain.SystemInstance{
		Ctx:                    context.Background(),
		DB:                     db,
		BookCategoryCacheStore: statsCacheStore,
		BookCategoryRepo:       bookCategoryRepo,
	}

	return testInstance
}

func prepareTestDB() {
	testInstance.DB.Delete(&model.BookCategory{})
}

func rollbackTestDb() {
	testInstance.DB.Close()
	testInstance = nil
}

var (
	bookCategoryCache = datarepo.UniqueKeyCacheDefinition{
		KeyPrefix:    "bc:",
		KeyFieldName: "ID",
		Expiration:   5 * time.Minute,
	}
)
