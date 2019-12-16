package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-txdb"
	redis2 "github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/merlinapp/datarepo-go"
	"github.com/merlinapp/datarepo-go/cachestore/redis"
	"github.com/merlinapp/datarepo-go/cachestore/stats"
	"github.com/merlinapp/datarepo-go/integration_tests/model"
	gorm2 "github.com/merlinapp/datarepo-go/repo/gorm"
	stats2 "github.com/merlinapp/datarepo-go/repo/stats"
	"log"
	"os"
	"strconv"
	"time"
)

type SystemInstance struct {
	Ctx                     context.Context
	DB                      *gorm.DB
	CacheStore              stats.StatsCacheStore
	UniqueKeyDataFetcher    stats2.StatsDataFetcher
	NonUniqueKeyDataFetcher stats2.StatsDataFetcher
	CachedRepo              datarepo.CachedRepository
}

var (
	testInstance      *SystemInstance
	connectionFactory ConnectionFactory
)

type ConnectionFactory func() *gorm.DB

func init() {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	connectionUrl := dbParams(dbUsername, dbPassword, dbHost, dbPort, dbName)
	log.Println("Connecting to DB: ", dbHost+":"+dbPort)

	noTxDbParam := os.Getenv("NO_TXDB")
	if noTxDbParam == "true" {
		connectionFactory = func() *gorm.DB {
			return ConnectLocalDb(connectionUrl)
		}
	} else {
		log.Println("Using TXDB - performing initial DB Automigrate using regular DB connection")
		db := ConnectLocalDb(connectionUrl)
		autoMigrate(db)
		db.Close()
		log.Println("DB Automigrate completed and connection closed")

		txdb.Register("txdb", "mysql", connectionUrl)
		connectionFactory = ConnectTxDb
	}
}

func StartSystemForIntegrationTests() *SystemInstance {
	if testInstance != nil {
		return testInstance
	}

	db := connectionFactory()
	db.LogMode(true)
	autoMigrate(db)

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

	testInstance = &SystemInstance{
		Ctx:                     context.Background(),
		DB:                      db,
		CacheStore:              statsCacheStore,
		CachedRepo:              repo,
		UniqueKeyDataFetcher:    statsUniqueKeyDataFetcher,
		NonUniqueKeyDataFetcher: statsNonUniqueKeyDataFetcher,
	}

	return testInstance
}

func PrepareTestDB() {
	testInstance.DB.Delete(&model.Book{})
}

func RollbackTestDb() {
	testInstance.DB.Close()
	testInstance = nil
}

func ConnectTxDb() *gorm.DB {
	sqlDb, err := sql.Open("txdb", fmt.Sprintf("connection_%d", time.Now().UnixNano()))
	if err != nil {
		log.Println(err)
		panic("failed to connect database")
	}

	db, err := gorm.Open("mysql", sqlDb)

	if err != nil {
		log.Println(err)
		panic("failed to connect database")
	}
	return db
}

func ConnectLocalDb(connectionUrl string) *gorm.DB {
	db, err := gorm.Open("mysql", connectionUrl)

	if err != nil {
		log.Println(err)
		panic("failed to connect database")
	}
	return db
}

func dbParams(params ...string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", params[0], params[1], params[2], params[3], params[4])
}

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(&model.Author{})
	db.AutoMigrate(&model.Book{})
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
