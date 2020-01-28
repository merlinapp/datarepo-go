package book_gorm_redis

import (
	"context"
	redis2 "github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/merlinapp/datarepo-go"
	"github.com/merlinapp/datarepo-go/cachestore/redis"
	"github.com/merlinapp/datarepo-go/cachestore/stats"
	"github.com/merlinapp/datarepo-go/integration_tests"
	"github.com/merlinapp/datarepo-go/integration_tests/book_gorm_redis/testdomain"
	"github.com/merlinapp/datarepo-go/integration_tests/model"
	gorm2 "github.com/merlinapp/datarepo-go/repo/gorm"
	stats2 "github.com/merlinapp/datarepo-go/repo/stats"
	"log"
	"os"
	"strconv"
	"time"
)

var testInstance *testdomain.SystemInstance

func startSystemForIntegrationTests() *testdomain.SystemInstance {
	if testInstance != nil {
		return testInstance
	}

	db := integration_tests.TestConnectionFactory()
	db.LogMode(true)

	database, err := strconv.Atoi(os.Getenv("REDIS_DATABASE"))
	if err != nil {
		log.Println("Couldn't parse the redis database number: ", os.Getenv("REDIS_DATABASE"))
		panic(err)
	}

	redisClient := redis2.NewClient(&redis2.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       database,
	})
	cacheStore := redis.NewRedisCacheStore(redisClient)
	statsCacheStore := stats.NewStatsCacheStore(cacheStore)

	gormUniqueKeyDataFetcher := gorm2.NewUniqueKeyDataFetcher(db, &model.Book{})
	gormNonUniqueKeyDataFetcher := gorm2.NewNonUniqueKeyDataFetcher(db, &model.Book{})
	gormDataWriter := gorm2.NewDataWriter(db, &model.Book{})

	statsUniqueKeyDataFetcher := stats2.NewStatsDataFetcher(gormUniqueKeyDataFetcher)
	statsNonUniqueKeyDataFetcher := stats2.NewStatsDataFetcher(gormNonUniqueKeyDataFetcher)
	statsDataWriter := stats2.NewStatsDataWriter(gormDataWriter)

	builder := datarepo.CachedRepositoryBuilder(&model.Book{}).
		WithUniqueKeyDataFetcher(statsUniqueKeyDataFetcher).
		WithNonUniqueKeyDataFetcher(statsNonUniqueKeyDataFetcher).
		WithDataWriter(statsDataWriter).
		WithUniqueKeyCache(idCache, statsCacheStore).
		WithNonUniqueKeyCache(authorIdCache, statsCacheStore)
	repo := builder.BuildCachedRepository()

	testInstance = &testdomain.SystemInstance{
		Ctx:                     context.Background(),
		DB:                      db,
		BookCacheStore:          statsCacheStore,
		BookRepo:                repo,
		UniqueKeyDataFetcher:    statsUniqueKeyDataFetcher,
		NonUniqueKeyDataFetcher: statsNonUniqueKeyDataFetcher,
	}

	return testInstance
}

func prepareTestDB() {
	testInstance.DB.Delete(&model.Book{})
}

func rollbackTestDb() {
	testInstance.DB.Close()
	testInstance = nil
}

var (
	idCache = datarepo.UniqueKeyCacheDefinition{
		KeyPrefix:    "b:",
		KeyFieldName: "ID",
		Expiration:   12 * time.Hour,
	}
	authorIdCache = datarepo.NonUniqueKeyCacheDefinition{
		KeyPrefix:         "a:",
		KeyFieldName:      "AuthorID",
		SubKeyFieldName:   "ID",
		Expiration:        12 * time.Hour,
		CacheEmptyResults: true,
	}
)
