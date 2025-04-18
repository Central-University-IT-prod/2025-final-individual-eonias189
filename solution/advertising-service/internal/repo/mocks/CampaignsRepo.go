// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	dto "advertising/advertising-service/internal/dto"
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "advertising/advertising-service/internal/models"

	uuid "github.com/google/uuid"
)

// CampaignsRepo is an autogenerated mock type for the CampaignsRepo type
type CampaignsRepo struct {
	mock.Mock
}

// CreateCampaign provides a mock function with given fields: ctx, advertiserId, data
func (_m *CampaignsRepo) CreateCampaign(ctx context.Context, advertiserId uuid.UUID, data dto.CampaignData) (uuid.UUID, error) {
	ret := _m.Called(ctx, advertiserId, data)

	if len(ret) == 0 {
		panic("no return value specified for CreateCampaign")
	}

	var r0 uuid.UUID
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, dto.CampaignData) (uuid.UUID, error)); ok {
		return rf(ctx, advertiserId, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, dto.CampaignData) uuid.UUID); ok {
		r0 = rf(ctx, advertiserId, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, dto.CampaignData) error); ok {
		r1 = rf(ctx, advertiserId, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteCampaign provides a mock function with given fields: ctx, campaignId
func (_m *CampaignsRepo) DeleteCampaign(ctx context.Context, campaignId uuid.UUID) error {
	ret := _m.Called(ctx, campaignId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteCampaign")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, campaignId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCampaignById provides a mock function with given fields: ctx, campaignId
func (_m *CampaignsRepo) GetCampaignById(ctx context.Context, campaignId uuid.UUID) (models.Campaign, error) {
	ret := _m.Called(ctx, campaignId)

	if len(ret) == 0 {
		panic("no return value specified for GetCampaignById")
	}

	var r0 models.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (models.Campaign, error)); ok {
		return rf(ctx, campaignId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) models.Campaign); ok {
		r0 = rf(ctx, campaignId)
	} else {
		r0 = ret.Get(0).(models.Campaign)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, campaignId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListCampaignsForAdvertiser provides a mock function with given fields: ctx, advertiserId, params
func (_m *CampaignsRepo) ListCampaignsForAdvertiser(ctx context.Context, advertiserId uuid.UUID, params dto.PaginationParams) ([]models.Campaign, error) {
	ret := _m.Called(ctx, advertiserId, params)

	if len(ret) == 0 {
		panic("no return value specified for ListCampaignsForAdvertiser")
	}

	var r0 []models.Campaign
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, dto.PaginationParams) ([]models.Campaign, error)); ok {
		return rf(ctx, advertiserId, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, dto.PaginationParams) []models.Campaign); ok {
		r0 = rf(ctx, advertiserId, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Campaign)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, dto.PaginationParams) error); ok {
		r1 = rf(ctx, advertiserId, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetCampaignAdImageUrl provides a mock function with given fields: ctx, campaignId, adImageUrl
func (_m *CampaignsRepo) SetCampaignAdImageUrl(ctx context.Context, campaignId uuid.UUID, adImageUrl *string) error {
	ret := _m.Called(ctx, campaignId, adImageUrl)

	if len(ret) == 0 {
		panic("no return value specified for SetCampaignAdImageUrl")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, *string) error); ok {
		r0 = rf(ctx, campaignId, adImageUrl)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateCampaign provides a mock function with given fields: ctx, campaignId, data
func (_m *CampaignsRepo) UpdateCampaign(ctx context.Context, campaignId uuid.UUID, data dto.CampaignData) error {
	ret := _m.Called(ctx, campaignId, data)

	if len(ret) == 0 {
		panic("no return value specified for UpdateCampaign")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, dto.CampaignData) error); ok {
		r0 = rf(ctx, campaignId, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewCampaignsRepo creates a new instance of CampaignsRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCampaignsRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *CampaignsRepo {
	mock := &CampaignsRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
