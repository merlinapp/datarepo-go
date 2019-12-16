package main

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/merlinapp/datarepo-go"
	redis2 "github.com/merlinapp/datarepo-go/cachestore/redis"
	"github.com/merlinapp/datarepo-go/cachestore/stats"
	"github.com/merlinapp/datarepo-go/example/entity"
	gorm2 "github.com/merlinapp/datarepo-go/repo/gorm"
	"github.com/satori/uuid"
	"log"
	"time"
)

func main() {
	// open a database connection and create our test database
	db := ConnectDb()
	db.LogMode(true)
	db.AutoMigrate(&entity.Book{})

	// open a connection to redis
	redisClient := ConnectRedis()

	// define the cache store we want to use
	cacheStore := redis2.NewRedisCacheStore(redisClient)
	statsCacheStore := stats.NewStatsCacheStore(cacheStore)

	// create the Cached datarepo
	// in this case we're using GORM and redis as our cache store
	repo := gorm2.CachedRepositoryBuilder(db, &entity.Book{}).
		WithUniqueKeyCache(idCache, statsCacheStore).
		WithNonUniqueKeyCache(authorIdCache, statsCacheStore).
		BuildCachedRepository()

	ctx := context.Background()

	authorId := uuid.NewV4().String()
	book := entity.Book{
		ID:       uuid.NewV4().String(),
		AuthorID: authorId,
		Status:   "completed",
	}

	// create a new book in our repo
	err := repo.Create(ctx, &book)
	if err != nil {
		log.Println("Error creating the book:", err)
		return
	}

	result, err := repo.FindByKey(ctx, "ID", book.ID)
	if err != nil || result.IsEmpty() {
		log.Println("Book not found or an error occurred!", err)
		return
	}
	var resultBook entity.Book
	result.InjectResult(&resultBook)
	log.Println("Stored Book: ", resultBook)

	result, err = repo.FindByKey(ctx, "AuthorID", authorId)
	if err != nil || result.IsEmpty() {
		log.Println("Couldn't fetch books by author!", err)
		return
	}
	var authorBooks []*entity.Book
	result.InjectResult(&authorBooks)
	log.Println("Author Book: ", *authorBooks[0])

	bookIds := []string{book.ID, uuid.NewV4().String()}
	bookResults, err := repo.FindByKeys(ctx, "ID", bookIds)
	if err != nil {
		log.Println("Couldn't fetch books by id!", err)
		return
	}
	var books []*entity.Book
	datarepo.InjectResults(bookResults, &books)
	log.Println("First Book: ", books[0])
	log.Println("Second Book: ", books[1])

	authorIds := []string{book.AuthorID, uuid.NewV4().String()}
	authorBooksResults, err := repo.FindByKeys(ctx, "AuthorID", authorIds)
	if err != nil {
		log.Println("Couldn't fetch books by authors!", err)
		return
	}
	var authorsBooks [][]*entity.Book
	datarepo.InjectResults(authorBooksResults, &authorsBooks)
	log.Println("Number of Books for first author: ", len(authorsBooks[0]))
	log.Println("Number of Books for second author: ", len(authorsBooks[1]))
}

func ConnectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "192.168.99.100:6379",
		Password: "", // no password set
		DB:       9,
	})
}

// Connect with mysql database
func ConnectDb() *gorm.DB {
	db, err := gorm.Open("mysql", "test:test1234@tcp(localhost:3306)/test?charset=utf8&parseTime=True")

	if err != nil {
		log.Println(err)
		panic("failed to connect database")
	}
	return db
}

var (
	idCache = datarepo.UniqueKeyCacheDefinition{
		KeyPrefix:    "b:",
		KeyFieldName: "ID",
		Expiration:   5 * time.Minute,
	}
	authorIdCache = datarepo.NonUniqueKeyCacheDefinition{
		KeyPrefix:         "a:",
		KeyFieldName:      "AuthorID",
		SubKeyFieldName:   "ID",
		Expiration:        5 * time.Minute,
		CacheEmptyResults: true,
	}
)
