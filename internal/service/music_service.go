package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Seven11Eleven/music_library/internal/domain/models"
	log "github.com/sirupsen/logrus"
	"net/url"
	"strings"
	"time"
)

type musicService struct {
	musicRepository       models.MusicRepository
	dataEnrichmentService DataEnrichmentService
}

func ValidateMusicName(musicName string) error {
	if musicName == "" {
		log.Warn("Validation failed: music song name is empty")
		return errors.New("music song name is required")
	}
	return nil
}

func ValidateMusicFilters(musicFilters models.MusicFilters) error {
	if musicFilters.ReleaseDate != nil && musicFilters.ReleaseDate.After(time.Now()) {
		log.Warnf("Validation failed: release date %v is in the future", musicFilters.ReleaseDate)
		return fmt.Errorf("release date cannot be in the future")
	}
	if musicFilters.Link != nil {
		_, err := url.ParseRequestURI(*musicFilters.Link)
		if err != nil {
			log.Warnf("Validation failed: music link %s is invalid", *musicFilters.Link)
			return fmt.Errorf("music link is invalid")
		}
	}

	if musicFilters.SongName != nil && len(*musicFilters.SongName) > 255 {
		log.Warnf("Validation failed: song name %s is too long", *musicFilters.SongName)
		return fmt.Errorf("song name must be shorter than 255 characters")
	}
	if musicFilters.GroupName != nil && len(*musicFilters.GroupName) > 255 {
		log.Warnf("Validation failed: group name %s is too long", *musicFilters.GroupName)
		return fmt.Errorf("group name must be shorter than 255 characters")
	}

	return nil
}

func ValidatePagination(limit, offset int) error {
	if limit <= 0 {
		log.Warn("Validation failed: limit must be greater than zero")
		return fmt.Errorf("limit must be greater than zero")
	}
	if offset < 0 {
		log.Warn("Validation failed: offset must be greater or equal to zero")
		return fmt.Errorf("offset must be greater or equal to zero")
	}
	return nil
}

func ValidateMusicID(musicID string) error {
	if musicID == "" {
		log.Warn("Validation failed: music song id is empty")
		return errors.New("music song id is required")
	}
	return nil
}

func checkContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		log.Warn("Context was cancelled")
		return ctx.Err()
	default:
		return nil
	}
}

func (m musicService) SaveMusic(ctx context.Context, music *models.MusicQuery) (*models.Music, error) {
	music.SongName = strings.ToLower(music.SongName)
	music.GroupName = strings.ToLower(music.GroupName)

	log.Infof("Saving new music: %s by %s", music.SongName, music.GroupName)

	err := ValidateMusicName(music.SongName)
	if err != nil {
		log.Warnf("Validation failed: %v", err)
		return nil, err
	}
	existingMusic, err := m.musicRepository.GetMusic(ctx, music.SongName, music.GroupName)
	if err != nil {
		log.Errorf("Error checking if music exists: %v", err)
		return nil, err
	}

	if existingMusic != nil {
		log.Infof("Music already exists: %s by %s", music.SongName, music.GroupName)
		return existingMusic, nil
	}

	enrichedMusic, err := m.dataEnrichmentService.FetchEnrichedMusic(ctx, music.GroupName, music.SongName)
	if err != nil {
		log.Errorf("Error during data enrichment: %v", err)
		return nil, err
	}

	enrichedMusic.SongName = strings.ToLower(enrichedMusic.SongName)
	enrichedMusic.GroupName = strings.ToLower(enrichedMusic.GroupName)

	res, err := m.musicRepository.SaveMusic(ctx, enrichedMusic)
	if err != nil {
		log.Errorf("Error saving music: %v", err)
		return nil, err
	}

	log.Infof("Music saved successfully: %s", res.SongName)
	return res, nil
}

func (m musicService) GetMusicsByFilters(ctx context.Context, filters models.MusicFilters, page, pageSize int) ([]models.Music, error) {
	log.Infof("Fetching music list with filters: %+v", filters)

	err := ValidateMusicFilters(filters)
	if err != nil {
		log.Warnf("Validation failed: %v", err)
		return nil, err
	}

	err = ValidatePagination(page, pageSize)
	if err != nil {
		log.Warnf("Pagination validation failed: %v", err)
		return nil, err
	}

	err = checkContext(ctx)
	if err != nil {
		log.Warnf("Context error: %v", err)
		return nil, err
	}

	res, err := m.musicRepository.GetMusicsByFilters(ctx, filters, page, pageSize)
	if err != nil {
		log.Errorf("Error fetching music list: %v", err)
		return nil, err
	}
	log.Infof("Successfully fetched %d music records", len(res))
	return res, nil
}

func (m musicService) GetMusicTextWithPaginationByVerse(ctx context.Context, musicID string, limit, offset int) (*models.Music, error) {
	log.Infof("Fetching verses for music ID: %s with pagination limit %d, offset %d", musicID, limit, offset)

	err := ValidateMusicID(musicID)
	if err != nil {
		log.Warnf("Validation failed: %v", err)
		return nil, err
	}

	err = ValidatePagination(limit, offset)
	if err != nil {
		log.Warnf("Pagination validation failed: %v", err)
		return nil, err
	}

	err = checkContext(ctx)
	if err != nil {
		log.Warnf("Context error: %v", err)
		return nil, err
	}

	res, err := m.musicRepository.GetMusicTextWithPaginationByVerse(ctx, musicID, limit, offset)
	if err != nil {
		log.Errorf("Error fetching verses for music ID %s: %v", musicID, err)
		return nil, err
	}

	log.Infof("Successfully fetched verses for music ID %s", musicID)
	return res, nil
}

func (m musicService) DeleteMusic(ctx context.Context, musicID string) error {
	log.Infof("Deleting music with ID: %s", musicID)

	err := ValidateMusicID(musicID)
	if err != nil {
		log.Warnf("Validation failed: %v", err)
		return err
	}

	err = m.musicRepository.DeleteMusic(ctx, musicID)
	if err != nil {
		log.Errorf("Error deleting music with ID %s: %v", musicID, err)
		return err
	}

	log.Infof("Music with ID %s deleted successfully", musicID)
	return nil
}

func (m musicService) UpdateMusic(ctx context.Context, music models.Music) (models.Music, error) {
	log.Infof("Updating music with ID: %s", music.ID)

	err := ValidateMusicID(music.ID)
	if err != nil {
		log.Warnf("Validation failed: %v", err)
		return models.Music{}, err
	}

	res, err := m.musicRepository.UpdateMusic(ctx, music)
	if err != nil {
		log.Errorf("Error updating music with ID %s: %v", music.ID, err)
		return models.Music{}, err
	}

	log.Infof("Music with ID %s updated successfully", music.ID)
	return res, nil
}

func NewMusicService(
	musicRepository models.MusicRepository,
	dataEnrichmentService DataEnrichmentService,
) models.MusicService {
	log.Info("Creating new music service")
	return &musicService{
		musicRepository:       musicRepository,
		dataEnrichmentService: dataEnrichmentService,
	}
}
