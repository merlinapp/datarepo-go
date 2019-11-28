package datarepo

// Defines a new Builder used to create a new CachedRepository
type Builder interface {
	// Adds a new Unique Key cache to the builder.
	//
	// Unique Key caches are used when a key in the cache refers to a single element.
	// For example, when the cache key represents a Primary Key of a database table.
	WithUniqueKeyCache(cacheDefinition UniqueKeyCacheDefinition, store CacheStore) Builder
	// Adds a new Non-Unique Key cache to the builder.
	//
	// Non-Unique Key caches are used when a key in the cache can hold more than one element.
	// For example, when the cache key represents a Foreign Key of a database table.
	//
	// Non-Unique caches require a subKey to be defined to compare the multiple elements inside a single cache
	// key entry.
	WithNonUniqueKeyCache(cacheDefinition NonUniqueKeyCacheDefinition, store CacheStore) Builder
	// Data Fetcher to use when retrieving data by a field considered a unique key
	WithUniqueKeyDataFetcher(fetcher DataFetcher) Builder
	// Data Fetcher to use when retrieving data by a field considered a non-unique key
	WithNonUniqueKeyDataFetcher(fetcher DataFetcher) Builder
	// DataWriter to be used when new data needs to be stored in a repository
	WithDataWriter(writer DataWriter) Builder
	// Indicates if the cache entries should be evicted entries after data is written to the repository
	// or if data in the cache should be updated instead
	EvictAfterWrite(v bool) Builder
	// Creates a new CachedRepository
	//
	// This method panics if it can't create the CachedRepository
	BuildCachedRepository() CachedRepository
	// Creates a new ReadOnlyCachedRepository.
	//
	// This method panics if it can't create the ReadOnlyCachedRepository
	BuildROCachedRepository() ReadOnlyCachedRepository
}

func CachedRepositoryBuilder(dataType interface{}) Builder {
	if dataType == nil {
		panic("The data type provided must not be nil")
	}

	builder := repositoryBuilder{
		DataType:        dataType,
		UniqueCaches:    make(map[string]uniqueCacheConfiguration),
		NonUniqueCaches: make(map[string]nonUniqueCacheConfiguration),
	}
	return &builder
}

type repositoryBuilder struct {
	DataType                interface{}
	UniqueKeyDataFetcher    DataFetcher
	NonUniqueKeyDataFetcher DataFetcher
	UniqueCaches            map[string]uniqueCacheConfiguration
	NonUniqueCaches         map[string]nonUniqueCacheConfiguration
	DataWriter              DataWriter
	EvictOnWrite            bool
}

type uniqueCacheConfiguration struct {
	UniqueKeyCacheDefinition
	CacheStore
}

type nonUniqueCacheConfiguration struct {
	NonUniqueKeyCacheDefinition
	CacheStore
}

func (b *repositoryBuilder) WithUniqueKeyDataFetcher(fetcher DataFetcher) Builder {
	b.UniqueKeyDataFetcher = fetcher
	return b
}

func (b *repositoryBuilder) WithNonUniqueKeyDataFetcher(fetcher DataFetcher) Builder {
	b.NonUniqueKeyDataFetcher = fetcher
	return b
}

func (b *repositoryBuilder) WithDataWriter(writer DataWriter) Builder {
	b.DataWriter = writer
	return b
}

func (b *repositoryBuilder) EvictAfterWrite(v bool) Builder {
	b.EvictOnWrite = v
	return b
}

func (b *repositoryBuilder) WithUniqueKeyCache(cacheDefinition UniqueKeyCacheDefinition, store CacheStore) Builder {
	b.validateCacheName(cacheDefinition.KeyFieldName)

	cacheConfig := uniqueCacheConfiguration{
		cacheDefinition,
		store,
	}
	b.UniqueCaches[cacheDefinition.KeyFieldName] = cacheConfig
	return b
}

func (b *repositoryBuilder) WithNonUniqueKeyCache(cacheDefinition NonUniqueKeyCacheDefinition, store CacheStore) Builder {
	b.validateCacheName(cacheDefinition.KeyFieldName)

	cacheConfig := nonUniqueCacheConfiguration{
		cacheDefinition,
		store,
	}
	b.NonUniqueCaches[cacheDefinition.KeyFieldName] = cacheConfig
	return b
}

func (b *repositoryBuilder) BuildCachedRepository() CachedRepository {
	if b.DataWriter == nil {
		panic("a DataWriter needs to be provided when building a new read-write cached repository")
	}
	roRepo := b.buildReadOnlyRepository()
	repo := cachedRepository{
		readOnlyCachedRepository: *roRepo,
		writer:                   b.DataWriter,
	}
	if b.EvictOnWrite {
		repo.postWriteOp = repo.evictFromCaches
	} else {
		repo.postWriteOp = repo.setValueInCaches
	}
	return &repo
}

func (b *repositoryBuilder) BuildROCachedRepository() ReadOnlyCachedRepository {
	if len(b.UniqueCaches) > 0 && b.UniqueKeyDataFetcher == nil {
		panic("a UniqueKeyDataFetcher needs to be provided when building a new cached repository with unique key caches defined")
	}
	if len(b.NonUniqueCaches) > 0 && b.NonUniqueKeyDataFetcher == nil {
		panic("a NonUniqueKeyDataFetcher needs to be provided when building a new cached repository with non-unique key caches defined")
	}
	return b.buildReadOnlyRepository()
}

func (b *repositoryBuilder) buildReadOnlyRepository() *readOnlyCachedRepository {
	repo := readOnlyCachedRepository{caches: make(map[string]Cache)}
	for k, v := range b.UniqueCaches {
		cacheHandler := UniqueKeyCache(b.DataType, v.UniqueKeyCacheDefinition)
		repo.caches[k] = Cache{
			Handler:     cacheHandler,
			Store:       v.CacheStore,
			DataFetcher: b.UniqueKeyDataFetcher,
		}
	}
	for k, v := range b.NonUniqueCaches {
		cacheHandler := NonUniqueKeyCache(b.DataType, v.NonUniqueKeyCacheDefinition)
		fetcher := b.NonUniqueKeyDataFetcher
		if v.CacheEmptyResults {
			fetcher = &emptyResultDataFetcherWrapper{
				dataType: b.DataType,
				delegate: fetcher,
			}
		}
		repo.caches[k] = Cache{
			Handler:     cacheHandler,
			Store:       v.CacheStore,
			DataFetcher: fetcher,
		}
	}
	return &repo
}

func (b *repositoryBuilder) validateCacheName(cacheName string) {
	if _, ok := b.UniqueCaches[cacheName]; ok {
		panic("a cache has already been defined with the same field name: " + cacheName)
	}
	if _, ok := b.NonUniqueCaches[cacheName]; ok {
		panic("a cache has already been defined with the same field name: " + cacheName)
	}
}
