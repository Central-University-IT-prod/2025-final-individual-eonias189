// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	models "advertising/advertising-service/internal/models"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// AdsRepo is an autogenerated mock type for the AdsRepo type
type AdsRepo struct {
	mock.Mock
}

// GetAdForClient provides a mock function with given fields: ctx, client, currentDay
func (_m *AdsRepo) GetAdForClient(ctx context.Context, client models.Client, currentDay int) (models.Ad, error) {
	ret := _m.Called(ctx, client, currentDay)

	if len(ret) == 0 {
		panic("no return value specified for GetAdForClient")
	}

	var r0 models.Ad
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Client, int) (models.Ad, error)); ok {
		return rf(ctx, client, currentDay)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.Client, int) models.Ad); ok {
		r0 = rf(ctx, client, currentDay)
	} else {
		r0 = ret.Get(0).(models.Ad)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.Client, int) error); ok {
		r1 = rf(ctx, client, currentDay)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAdsRepo creates a new instance of AdsRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAdsRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *AdsRepo {
	mock := &AdsRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
