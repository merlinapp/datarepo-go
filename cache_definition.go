package datarepo

import "time"

type UniqueKeyCacheDefinition struct {
	KeyPrefix string
	// Name of the field that defines the key that will be used to store elements in the cache
	KeyFieldName string
	// expiration time of entries in the cache
	Expiration time.Duration
}

type NonUniqueKeyCacheDefinition struct {
	KeyPrefix string
	// Name of the field that defines the key that will be used to store elements in the cache
	KeyFieldName string
	// Name of the field that defines the subkey that will be used to compare and store elements that
	// belong to the same key
	SubKeyFieldName string
	// Expiration time of entries in the cache
	Expiration time.Duration
	// Indicates if empty results should be cached
	CacheEmptyResults bool
}
