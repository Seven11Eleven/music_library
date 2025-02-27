// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/Seven11Eleven/music_library/internal/domain/models"
	mock "github.com/stretchr/testify/mock"
)

// MusicRepository is an autogenerated mock type for the MusicRepository type
type MusicRepository struct {
	mock.Mock
}

// DeleteMusic provides a mock function with given fields: ctx, musicID
func (_m *MusicRepository) DeleteMusic(ctx context.Context, musicID string) error {
	ret := _m.Called(ctx, musicID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteMusic")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, musicID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetMusic provides a mock function with given fields: ctx, musicName, groupName
func (_m *MusicRepository) GetMusic(ctx context.Context, musicName string, groupName string) (*models.Music, error) {
	ret := _m.Called(ctx, musicName, groupName)

	if len(ret) == 0 {
		panic("no return value specified for GetMusic")
	}

	var r0 *models.Music
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*models.Music, error)); ok {
		return rf(ctx, musicName, groupName)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *models.Music); ok {
		r0 = rf(ctx, musicName, groupName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Music)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, musicName, groupName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMusicTextWithPaginationByVerse provides a mock function with given fields: ctx, musicID, limit, offset
func (_m *MusicRepository) GetMusicTextWithPaginationByVerse(ctx context.Context, musicID string, limit int, offset int) (*models.Music, error) {
	ret := _m.Called(ctx, musicID, limit, offset)

	if len(ret) == 0 {
		panic("no return value specified for GetMusicTextWithPaginationByVerse")
	}

	var r0 *models.Music
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) (*models.Music, error)); ok {
		return rf(ctx, musicID, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) *models.Music); ok {
		r0 = rf(ctx, musicID, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Music)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) error); ok {
		r1 = rf(ctx, musicID, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMusicsByFilters provides a mock function with given fields: ctx, filters, page, pageSize
func (_m *MusicRepository) GetMusicsByFilters(ctx context.Context, filters models.MusicFilters, page int, pageSize int) ([]models.Music, error) {
	ret := _m.Called(ctx, filters, page, pageSize)

	if len(ret) == 0 {
		panic("no return value specified for GetMusicsByFilters")
	}

	var r0 []models.Music
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.MusicFilters, int, int) ([]models.Music, error)); ok {
		return rf(ctx, filters, page, pageSize)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.MusicFilters, int, int) []models.Music); ok {
		r0 = rf(ctx, filters, page, pageSize)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Music)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.MusicFilters, int, int) error); ok {
		r1 = rf(ctx, filters, page, pageSize)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SaveMusic provides a mock function with given fields: ctx, music
func (_m *MusicRepository) SaveMusic(ctx context.Context, music *models.Music) (*models.Music, error) {
	ret := _m.Called(ctx, music)

	if len(ret) == 0 {
		panic("no return value specified for SaveMusic")
	}

	var r0 *models.Music
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Music) (*models.Music, error)); ok {
		return rf(ctx, music)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Music) *models.Music); ok {
		r0 = rf(ctx, music)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Music)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Music) error); ok {
		r1 = rf(ctx, music)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateMusic provides a mock function with given fields: ctx, music
func (_m *MusicRepository) UpdateMusic(ctx context.Context, music models.Music) (models.Music, error) {
	ret := _m.Called(ctx, music)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMusic")
	}

	var r0 models.Music
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.Music) (models.Music, error)); ok {
		return rf(ctx, music)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.Music) models.Music); ok {
		r0 = rf(ctx, music)
	} else {
		r0 = ret.Get(0).(models.Music)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.Music) error); ok {
		r1 = rf(ctx, music)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMusicRepository creates a new instance of MusicRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMusicRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MusicRepository {
	mock := &MusicRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
