package datarepo

import (
	"context"
)

type CachedRepository interface {
	ReadOnlyCachedRepository
	// Inserts the provided value into the repository
	//
	// If successful, then the value is also cached in the configured caches and/or evicted
	// from the caches that requested eviction on write operations
	Create(ctx context.Context, value interface{}) error
	// Updates the provided value in the repository
	//
	// If successful, then the value is also updated/cached in the configured caches and/or evicted
	// from the caches that requested eviction on write operations
	Update(ctx context.Context, value interface{}) error
	//Updates partially the provided value in the repository
	//
	// If successful, then the value is also updated/cached in the configured caches and/or evicted
	// Also fully updated object will be retrieved from data storage and will be cached
	// value param works as in/out: full object, retrieved from database after update, will be stored in this variable
	// from the caches that requested eviction on write operations
	PartialUpdate(ctx context.Context, value interface{}) error
}

type cachedRepository struct {
	readOnlyCachedRepository
	postWriteOp postWriteOperation
	writer      DataWriter
}

type postWriteOperation func(ctx context.Context, value interface{}) error

func (r *cachedRepository) Create(ctx context.Context, value interface{}) error {
	err := r.writer.Create(ctx, value)
	if err != nil {
		return err
	}

	err = r.postWriteOp(ctx, value)
	if err != nil {
		return err
	}
	return nil
}

func (r *cachedRepository) Update(ctx context.Context, value interface{}) error {
	err := r.writer.Update(ctx, value)
	if err != nil {
		return err
	}

	err = r.postWriteOp(ctx, value)
	if err != nil {
		return err
	}
	return nil
}

func (r *cachedRepository) PartialUpdate(ctx context.Context, value interface{}) error {
	err := r.writer.PartialUpdate(ctx, value)
	if err != nil {
		return err
	}

	err = r.postWriteOp(ctx, value)
	if err != nil {
		return err
	}

	return nil
}

func (r *cachedRepository) setValueInCaches(ctx context.Context, value interface{}) error {
	for _, v := range r.caches {
		err := v.Set(ctx, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *cachedRepository) evictFromCaches(ctx context.Context, value interface{}) error {
	for _, v := range r.caches {
		err := v.DeleteValue(ctx, value)
		if err != nil {
			return err
		}
	}
	return nil
}
