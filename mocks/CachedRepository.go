// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import datarepo "github.com/merlinapp/datarepo-go"
import mock "github.com/stretchr/testify/mock"

// CachedRepository is an autogenerated mock type for the CachedRepository type
type CachedRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, value
func (_m *CachedRepository) Create(ctx context.Context, value interface{}) error {
	ret := _m.Called(ctx, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindByKey provides a mock function with given fields: ctx, keyFieldName, id
func (_m *CachedRepository) FindByKey(ctx context.Context, keyFieldName string, id interface{}) (datarepo.Result, error) {
	ret := _m.Called(ctx, keyFieldName, id)

	var r0 datarepo.Result
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}) datarepo.Result); ok {
		r0 = rf(ctx, keyFieldName, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(datarepo.Result)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, interface{}) error); ok {
		r1 = rf(ctx, keyFieldName, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByKeys provides a mock function with given fields: ctx, keyFieldName, ids
func (_m *CachedRepository) FindByKeys(ctx context.Context, keyFieldName string, ids interface{}) ([]datarepo.Result, error) {
	ret := _m.Called(ctx, keyFieldName, ids)

	var r0 []datarepo.Result
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}) []datarepo.Result); ok {
		r0 = rf(ctx, keyFieldName, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]datarepo.Result)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, interface{}) error); ok {
		r1 = rf(ctx, keyFieldName, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PartialUpdate provides a mock function with given fields: ctx, value
func (_m *CachedRepository) PartialUpdate(ctx context.Context, value interface{}) error {
	ret := _m.Called(ctx, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, value
func (_m *CachedRepository) Update(ctx context.Context, value interface{}) error {
	ret := _m.Called(ctx, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) error); ok {
		r0 = rf(ctx, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
