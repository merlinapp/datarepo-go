# datarepo-go
A data access/repository library with caching support

# Why datarepo-go?

A use case in our data access layers is to have caching support and the implementations tend to be very similar, or should be very similar in their behavior. For example, for a read operation, we check in the cache and if the data is not there we fallback to the data repository (e.g. a MySQL database).

This logic of checking/adding data to a cache can be centralized and is exactly the intention of this library.

The library is based on configurable key-based access to the data to allow for a unified caching strategy. We'll cover this further in our examples section.

# Status

This is a new library and is considered to be in **Beta** and the APIs can change as we adapt more use cases.

We'll follow semantic versioning, so, once we reach a stable state and release 1.0.0 API's should be stable inside releases with the same major number.

# Features

Right now we offer support for:

Caching Stores:
* Redis Cache (using go-redis and implemented as a write-through cache)
* In-Memory Cache (using freecache - see github.com/coocood/freecache)
* Statistics Wrapper (a Caching Store that provides stats about cache access, useful for testing)

Repositories:
* GORM-based repo (any DB supported by GORM)
* Statistics Wrappers (wrappers around the repository access classes that provide stats about access to the repositories, useful for testing)

Cache Types:
* Unique Key caches, where a given key results in a single entity/instance. For example, the BookID in a Book entity.
* Non-Unique Key caches, where a given key can result in multiple instances. For example, the AuthorID in a Book entity, where a single author can have one or more books.

Implementation Note:
Non-Unique Key caches should be used with care and only where the cardinality isn't too large. For example, if one key can have a million records/instances associated
a different strategy should be considered.

# An example use case

The example code can be found in [example/main.go](./example/main.go)

Assume you have a Book entity that is queried very often, both by bookId and also by author. In our example use case, the number of books an author has doesn't change very often so it doesn't make sense to hit the database for every query.

```go
type Book struct {
	ID       string `json:"id" gorm:"primary_key" sql:"type:CHAR(36)"`
	AuthorID string `json:"authorId" sql:"type:CHAR(36)"`
	Status   string `json:"status""`
}
```

In our example use case we're using Gorm to store data in a MySQL database and Redis as the store for our caching layer.

First, we define our Gorm DB connection and Redis connections:

```go
// open a database connection and create our test database
db := ConnectDb()
db.LogMode(true)
db.AutoMigrate(&entity.Book{})

// open a connection to redis
redisClient := ConnectRedis()
```

Second, we define the cache store we want to use. In our case, we want to create a cache store backed by our redis client:
```go
cacheStore := redis2.NewRedisCacheStore(redisClient)
statsCacheStore := stats.NewStatsCacheStore(cacheStore)
```

The `statsCacheStore` is useful for testing purposes and it's just a wrapper/decorator around the main `cacheStore`.

In the code below, you could replace `statsCacheStore` with `cacheStore` if you don't want any cache stats.

Third, we define the caches we want to keep:

```go
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
```

And finally, we create our cached data repository to handle `Book` instances:

```go
repo := gorm.CachedRepositoryBuilder(db, &entity.Book{}).
    WithUniqueKeyCache(idCache, statsCacheStore).
    WithNonUniqueKeyCache(authorIdCache, statsCacheStore).
    BuildCachedRepository()
```

In this case, we're defining two caches that we want to keep:

## Unique Key Cache

We have a Unique Key Cache called `idCache`. Unique Key Caches are those that for every key in the cache we expect a single entry/instance.

In our case, we're using the book `ID` as the field we should use for the cache key. As this is the primary key of the entity we can safely use it
as a Unique Key Cache.

An example key in Redis would be: `b:cddb0298-7d55-4e96-be32-2cbfa30ec12d` and it is guaranteed to contain a single book.

## Non-Unique Key Caches

We also have a Non-Unique Key Cache called `authorIdCache`. Non-Unique Key Caches are those that for every key in the cache we can expect one or more entries.

In our case, we're using the book's `AuthorID` field as the cache key. Of course, one author can have multiple books. 

An example key in Redis would be: `a:804b2cdb-5a4e-4845-a9a5-a097ac2322ac` and it is guaranteed to contain an array of books. This array can be empty if we enable the `CacheEmptyResults` flag.

Non-Unique Key Caches need to define a `SubKeyFieldName` that is used to compare the books inside the array. This should in general be the Primary Key or a unique key of the entity. In our case, our `SubKeyFieldName` is set to be the `ID` field of our book.

# Using the repo

Once you have your cached data repository you'll have an instance that implements the `CachedDataRepository` interface which currently provides the following 4 methods:

```go
Create(ctx context.Context, value interface{}) error
Update(ctx context.Context, value interface{}) error
FindByKey(ctx context.Context, keyFieldName string, id interface{}) (Result, error)
FindByKeys(ctx context.Context, keyFieldName string, ids interface{}) ([]Result, error)
```

## Create/Update data

If you're familiar with Gorm, this should feel familiar to some extent. The main difference is that currently we only support updating an object fully.

For example, to create a book:
```go
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
```

Once data is created/updated the configured caches will be updated to include the new record.

If on the other hand you want to evict the data from the caches when data is written, you can use the following line in your repository builder:

```go
// by default this is false
builder.EvictAfterWrite(true)
```

### Type checking when writing data

If you pass an element other than a `*Book` to the `Create` and `Update` methods the Gorm repository will return an error. It does a type check to ensure that the element you're trying to store is of the expected type the repo was created with. 

### What caches are updated when `EvictAfterWrite` is false?

Unique Key Caches are guaranteed to be updated on every write.

Non-Unique Key Caches are updated only if the key already exists in the cache. 

For example, if the cache entry for an author already exists, let's say: `a:x` -> `[ book1, book2 ]`, adding a new book to the same author will correctly append the 3rd book to the cached array: `a:x` -> `[book1, book2, book3]`. But, if the cache key for that author (`a:x`) doesn't exist then no data will be cached. 

## Reading data

You can read data of a single or multiple ids using the `FindByKey` and `FindByKeys` methods respectively.

## Single Key from Unique Key Cache : Fetching a single book by ID

For example, to read the data of a single book by id:

```go
result, err := repo.FindByKey(ctx, "ID", book.ID)
// error handling goes here...
var resultBook entity.Book
result.InjectResult(&resultBook)
log.Println("Stored Book: ", resultBook)
```

NOTE: The second parameter to `FindByKey` must correspond to a `Book` field that has a cache defined. In our case, `ID` has our `idCache` defined in our repository.

NOTE: The method `InjectResult` is used to transform a single `Result` into an element of the expected type, in our case `*entity.Book`. Make sure you check that the result isn't empty before invoking this method.

If you run the example code, you'll see that there is no database access in this case as we just created the book and is available in the cache.

## Multiple Keys from Unique Key Cache : Fetching multiple books by ID

If you want to fetch multiple books by id:

```go
bookIds := []string{book.ID, uuid.NewV4().String()}
bookResults, err := repo.FindByKeys(ctx, "ID", bookIds)
// error handling goes here...
var books []*entity.Book
datarepo.InjectResults(bookResults, &books)
log.Println("First Book: ", books[0])
log.Println("Second Book: ", books[1])
```

NOTE: `datarepo.InjectResults` is a utility method to transform an slice of `Result` instances into instances of the expected data type, in our case a slice `[]*entity.Book`.

If you run the example code you'll see that the first book is retrieved from the cache. The second book wasn't found in the cache so we try to fetch it from the database and you'll see the respective SQL statement from Gorm.

Result orders matter, in all calls to `FindByKeys`, the `i-th` item in the returned slice corresponds to the `i-th` item in the input `ids` slice.

## Single Key from Non-Unique Key Cache : Fetching books of a single author by ID

To fetch the books of a single author:

```go
result, err = repo.FindByKey(ctx, "AuthorID", authorId)
// error handling goes here...
var authorBooks []*entity.Book
result.InjectResult(&authorBooks)
log.Println("Author Book: ", *authorBooks[0])
```

In this case, the field we want to use to perform the query is `AuthorID` which is backed by the `authorIdCache` non-unique key cache.

## Multiple Key from Non-Unique Key Cache: Fetching books of multiple authors by ID

To fetch the books of multiple authors:

```go
authorIds := []string{book.AuthorID, uuid.NewV4().String()}
authorBooksResults, err := repo.FindByKeys(ctx, "AuthorID", authorIds)
// error handling goes here...
var authorsBooks [][]*entity.Book
datarepo.InjectResults(authorBooksResults, &authorsBooks)
log.Println("Number of Books for first author: ", len(authorsBooks[0]))
log.Println("Number of Books for second author: ", len(authorsBooks[1]))
```

In this case our expected result is a slice of slices (`[][]*entity.Book`). For each author id we provide, we retrieve a slice of books for that author.
