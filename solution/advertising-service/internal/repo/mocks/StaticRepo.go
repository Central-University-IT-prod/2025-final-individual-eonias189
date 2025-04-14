// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	models "advertising/advertising-service/internal/models"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// StaticRepo is an autogenerated mock type for the StaticRepo type
type StaticRepo struct {
	mock.Mock
}

// DeleteStatic provides a mock function with given fields: ctx, name
func (_m *StaticRepo) DeleteStatic(ctx context.Context, name string) error {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for DeleteStatic")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LoadStatic provides a mock function with given fields: ctx, name
func (_m *StaticRepo) LoadStatic(ctx context.Context, name string) (models.Static, error) {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for LoadStatic")
	}

	var r0 models.Static
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (models.Static, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) models.Static); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(models.Static)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveStatic provides a mock function with given fields: ctx, name, static
func (_m *StaticRepo) SaveStatic(ctx context.Context, name string, static models.Static) error {
	ret := _m.Called(ctx, name, static)

	if len(ret) == 0 {
		panic("no return value specified for SaveStatic")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, models.Static) error); ok {
		r0 = rf(ctx, name, static)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStaticRepo creates a new instance of StaticRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStaticRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *StaticRepo {
	mock := &StaticRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
