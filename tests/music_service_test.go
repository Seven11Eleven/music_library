package service_test

import (
	"context"
	"github.com/Seven11Eleven/music_library/internal/domain/mocks"
	"github.com/Seven11Eleven/music_library/internal/domain/models"
	"github.com/Seven11Eleven/music_library/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSaveMusic(t *testing.T) {
	musicRepoMock := new(mocks.MusicRepository)
	dataEnrichmentMock := new(mocks.DataEnrichmentService)

	existingMusic := &models.Music{
		SongName:  "sonne",
		GroupName: "rammstein",
	}

	musicRepoMock.On("GetMusic", mock.Anything, "sonne", "rammstein").Return(existingMusic, nil)

	musicService := service.NewMusicService(musicRepoMock, dataEnrichmentMock)

	newMusic := &models.MusicQuery{
		SongName:  "sonne",
		GroupName: "rammstein",
	}

	result, err := musicService.SaveMusic(context.TODO(), newMusic)

	assert.NoError(t, err)
	assert.Equal(t, existingMusic, result)

	dataEnrichmentMock.AssertNotCalled(t, "FetchEnrichedMusic", mock.Anything, "rammstein", "sonne")

	musicRepoMock.AssertExpectations(t)
}

func TestSaveMusic_ValidationFailed(t *testing.T) {
	ctx := context.TODO()

	mockMusicRepo := new(mocks.MusicRepository)
	mockDataEnrichmentService := new(mocks.DataEnrichmentService)

	musicService := service.NewMusicService(mockMusicRepo, mockDataEnrichmentService)

	musicQuery := &models.MusicQuery{
		SongName:  "",
		GroupName: "Rammstein",
	}

	_, err := musicService.SaveMusic(ctx, musicQuery)

	assert.EqualError(t, err, "music song name is required")

	mockMusicRepo.AssertNotCalled(t, "GetMusic", mock.Anything, mock.Anything, mock.Anything)
	mockDataEnrichmentService.AssertNotCalled(t, "FetchEnrichedMusic", mock.Anything, mock.Anything, mock.Anything)
}

func TestGetMusicsByFilters(t *testing.T) {
	ctx := context.TODO()
	mockMusicRepo := new(mocks.MusicRepository)
	mockDataEnrichmentService := new(mocks.DataEnrichmentService)

	mockMusicRepo.On("GetMusicsByFilters", ctx, mock.Anything, 1, 10).
		Return([]models.Music{
			{ID: "1", SongName: "Sonne", GroupName: "Rammstein"},
			{ID: "2", SongName: "Du Hast", GroupName: "Rammstein"},
		}, nil)

	musicService := service.NewMusicService(mockMusicRepo, mockDataEnrichmentService)

	filters := models.MusicFilters{}
	musics, err := musicService.GetMusicsByFilters(ctx, filters, 1, 10)

	assert.NoError(t, err)
	assert.Len(t, musics, 2)
	assert.Equal(t, "Sonne", musics[0].SongName)
	assert.Equal(t, "Du Hast", musics[1].SongName)

	mockMusicRepo.AssertExpectations(t)
}

func TestGetMusicTextWithPaginationByVerse(t *testing.T) {
	ctx := context.TODO()
	mockMusicRepo := new(mocks.MusicRepository)
	mockDataEnrichmentService := new(mocks.DataEnrichmentService)

	// Мокаем получение стихов
	mockMusicRepo.On("GetMusicTextWithPaginationByVerse", ctx, "1", 10, 0).
		Return(&models.Music{
			ID:       "1",
			SongName: "Sonne",
			Verses: []models.Verse{
				{Text: "Eins, zwei, drei", Number: 1},
			},
		}, nil)

	// Создаем сервис с моками
	musicService := service.NewMusicService(mockMusicRepo, mockDataEnrichmentService)

	// Валидный запрос для получения стихов
	music, err := musicService.GetMusicTextWithPaginationByVerse(ctx, "1", 10, 0)

	assert.NoError(t, err)
	assert.Equal(t, "Sonne", music.SongName)
	assert.Len(t, music.Verses, 1)
	assert.Equal(t, "Eins, zwei, drei", music.Verses[0].Text)

	// Проверяем моки
	mockMusicRepo.AssertExpectations(t)
}

func TestDeleteMusic(t *testing.T) {
	ctx := context.TODO()
	mockMusicRepo := new(mocks.MusicRepository)
	mockDataEnrichmentService := new(mocks.DataEnrichmentService)

	mockMusicRepo.On("DeleteMusic", ctx, "1").Return(nil)

	musicService := service.NewMusicService(mockMusicRepo, mockDataEnrichmentService)
	err := musicService.DeleteMusic(ctx, "1")

	assert.NoError(t, err)

	mockMusicRepo.AssertExpectations(t)
}

func TestUpdateMusic(t *testing.T) {
	ctx := context.TODO()
	mockMusicRepo := new(mocks.MusicRepository)
	mockDataEnrichmentService := new(mocks.DataEnrichmentService)

	mockMusicRepo.On("UpdateMusic", ctx, mock.Anything).
		Return(models.Music{ID: "1", SongName: "Sonne", GroupName: "Rammstein"}, nil)

	musicService := service.NewMusicService(mockMusicRepo, mockDataEnrichmentService)

	music := models.Music{ID: "1", SongName: "Sonne", GroupName: "Rammstein"}
	updatedMusic, err := musicService.UpdateMusic(ctx, music)

	assert.NoError(t, err)
	assert.Equal(t, "Sonne", updatedMusic.SongName)

	mockMusicRepo.AssertExpectations(t)
}
