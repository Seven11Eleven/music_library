// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/Seven11Eleven/music_library/internal/domain/models"
	mock "github.com/stretchr/testify/mock"
)

// DataEnrichmentService is an autogenerated mock type for the DataEnrichmentService type
type DataEnrichmentService struct {
	mock.Mock
}

// FetchEnrichedMusic provides a mock function with given fields: ctx, groupName, musicName
func (_m *DataEnrichmentService) FetchEnrichedMusic(ctx context.Context, groupName string, musicName string) (*models.Music, error) {
	ret := _m.Called(ctx, groupName, musicName)

	if len(ret) == 0 {
		panic("no return value specified for FetchEnrichedMusic")
	}

	var r0 *models.Music
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*models.Music, error)); ok {
		return rf(ctx, groupName, musicName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *models.Music); ok {
		r0 = rf(ctx, groupName, musicName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Music)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, groupName, musicName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewDataEnrichmentService creates a new instance of DataEnrichmentService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDataEnrichmentService(t interface {
	mock.TestingT
	Cleanup(func())
}) *DataEnrichmentService {
	mock := &DataEnrichmentService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
