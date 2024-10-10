package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Seven11Eleven/music_library/internal/config"
	"github.com/Seven11Eleven/music_library/internal/domain/models"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type DataEnrichmentService interface {
	FetchEnrichedMusic(ctx context.Context, groupName, musicName string) (*models.Music, error)
}

type dataEnrichmentService struct {
	config *config.Config
}

func (d dataEnrichmentService) FetchEnrichedMusic(_ context.Context, groupName, musicName string) (*models.Music, error) {
	log.Infof("Starting enrichment for song '%s' by group '%s'", musicName, groupName)

	encodedGroupName := url.QueryEscape(groupName)
	encodedMusicName := url.QueryEscape(musicName)

	apiKey := d.config.APIKey
	urlDetails := fmt.Sprintf("http://ws.audioscrobbler.com/2.0/?method=track.getInfo&api_key=%s&artist=%s&track=%s&format=json", apiKey, encodedGroupName, encodedMusicName)
	urlLyrics := fmt.Sprintf("https://lyrist.vercel.app/api/%s/%s", encodedMusicName, encodedGroupName)

	log.Infof("Fetching track details from URL: %s", urlDetails)
	resDetails, err := http.Get(urlDetails)
	if err != nil {
		log.Errorf("Failed to fetch track details: %v", err)
		return nil, err
	}
	defer resDetails.Body.Close()

	if resDetails.StatusCode != http.StatusOK {
		log.Errorf("Received non-200 status code from track details API: %s", resDetails.Status)
		return nil, fmt.Errorf(resDetails.Status)
	}

	var trackData struct {
		Track struct {
			URL  string `json:"url"`
			Wiki struct {
				Published string `json:"published"`
			} `json:"wiki"`
		} `json:"track"`
	}

	log.Info("Parsing track details response...")
	if err := json.NewDecoder(resDetails.Body).Decode(&trackData); err != nil {
		log.Errorf("Failed to parse track details: %v", err)
		return nil, err
	}

	if trackData.Track.URL == "" {
		log.Warnf("No track URL found for song '%s' by group '%s'", musicName, groupName)
		return nil, fmt.Errorf("no track URL found for song %s by group %s", musicName, groupName)
	}

	var releaseDate *time.Time
	if trackData.Track.Wiki.Published != "" {
		log.Infof("Parsing release date: %s", trackData.Track.Wiki.Published)
		parsedDate, err := time.Parse("2 Jan 2006, 15:04", trackData.Track.Wiki.Published)
		if err != nil {
			log.Warnf("Error parsing release date, using default: %v", err)
		} else {
			releaseDate = &parsedDate
		}
	}

	log.Infof("Fetching lyrics from URL: %s", urlLyrics)
	resLyrics, err := http.Get(urlLyrics)
	if err != nil {
		log.Errorf("Failed to fetch lyrics: %v", err)
		return nil, err
	}
	defer resLyrics.Body.Close()

	if resLyrics.StatusCode != http.StatusOK {
		log.Errorf("Received non-200 response code from lyrics API")
		return nil, fmt.Errorf("received non-200 response code from lyrics API")
	}

	var lyricsData struct {
		Lyrics string `json:"lyrics"`
	}

	log.Info("Parsing lyrics response...")
	if err := json.NewDecoder(resLyrics.Body).Decode(&lyricsData); err != nil {
		log.Errorf("Failed to parse lyrics: %v", err)
		return nil, err
	}

	if strings.TrimSpace(lyricsData.Lyrics) == "" {
		log.Warnf("No lyrics found for song '%s' by group '%s'", musicName, groupName)
		return nil, fmt.Errorf("no lyrics found for song %s by group %s", musicName, groupName)
	}

	log.Info("Parsing verses from lyrics...")
	verses := parseVerses(lyricsData.Lyrics)

	music := &models.Music{
		ReleaseDate: releaseDate,
		Verses:      verses,
		Link:        trackData.Track.URL,
		SongName:    musicName,
		GroupName:   groupName,
	}

	log.Infof("Successfully enriched music for song '%s' by group '%s'", musicName, groupName)
	return music, nil
}

func parseVerses(songText string) []models.Verse {
	verseTexts := strings.Split(songText, "\n\n")
	verses := make([]models.Verse, 0, len(verseTexts))

	for i, verseText := range verseTexts {
		verse := models.Verse{
			Text:   verseText,
			Number: i,
		}
		verses = append(verses, verse)
	}

	log.Infof("Parsed %d verses from song text", len(verses))
	return verses
}

func NewDataEnrichmentService(cfg *config.Config) DataEnrichmentService {
	log.Info("Creating new data enrichment service")
	return &dataEnrichmentService{config: cfg}
}
