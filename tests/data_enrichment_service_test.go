package service_test

import (
	"context"
	"github.com/Seven11Eleven/music_library/internal/config"
	"github.com/Seven11Eleven/music_library/internal/service"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestFetchEnrichedMusic_Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := &config.Config{
		APIKey: "test_api_key",
	}

	dataEnrichmentService := service.NewDataEnrichmentService(cfg)

	httpmock.RegisterResponder("GET", "http://ws.audioscrobbler.com/2.0/?method=track.getInfo&api_key=test_api_key&artist=Rammstein&track=Sonne&format=json",
		httpmock.NewStringResponder(http.StatusOK, `{
			"track": {
				"url": "http://www.example.com",
				"wiki": {
					"published": "2 Jan 2006, 15:04"
				}
			}
		}`))

	httpmock.RegisterResponder("GET", "https://lyrist.vercel.app/api/Sonne/Rammstein",
		httpmock.NewStringResponder(http.StatusOK, `{
			"lyrics": "[Verse 1]\nThis is the first verse\n\n[Verse 2]\nThis is the second verse"
		}`))

	music, err := dataEnrichmentService.FetchEnrichedMusic(context.Background(), "Rammstein", "Sonne")
	assert.NoError(t, err)
	assert.NotNil(t, music)
	assert.Equal(t, "http://www.example.com", music.Link)
	assert.Len(t, music.Verses, 2)

	assert.Equal(t, "[Verse 1]\nThis is the first verse", music.Verses[0].Text)
	assert.Equal(t, "[Verse 2]\nThis is the second verse", music.Verses[1].Text)

	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info["GET http://ws.audioscrobbler.com/2.0/?method=track.getInfo&api_key=test_api_key&artist=Rammstein&track=Sonne&format=json"])
	assert.Equal(t, 1, info["GET https://lyrist.vercel.app/api/Sonne/Rammstein"])
}

func TestFetchEnrichedMusic_TrackAPIError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := &config.Config{
		APIKey: "test_api_key",
	}

	dataEnrichmentService := service.NewDataEnrichmentService(cfg)

	httpmock.RegisterResponder("GET", "http://ws.audioscrobbler.com/2.0/?method=track.getInfo&api_key=test_api_key&artist=Rammstein&track=Sonne&format=json",
		httpmock.NewStringResponder(http.StatusInternalServerError, ""))

	music, err := dataEnrichmentService.FetchEnrichedMusic(context.Background(), "Rammstein", "Sonne")
	assert.Error(t, err)
	assert.Nil(t, music)

	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info["GET http://ws.audioscrobbler.com/2.0/?method=track.getInfo&api_key=test_api_key&artist=Rammstein&track=Sonne&format=json"])
	assert.Equal(t, 0, info["GET https://lyrist.vercel.app/api/Sonne/Rammstein"])
}

func TestFetchEnrichedMusic_LyricsAPIError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := &config.Config{
		APIKey: "test_api_key",
	}

	dataEnrichmentService := service.NewDataEnrichmentService(cfg)

	httpmock.RegisterResponder("GET", "http://ws.audioscrobbler.com/2.0/?method=track.getInfo&api_key=test_api_key&artist=Rammstein&track=Sonne&format=json",
		httpmock.NewStringResponder(http.StatusOK, `{
			"track": {
				"url": "http://www.example.com",
				"wiki": {
					"published": "2 Jan 2006, 15:04"
				}
			}
		}`))

	httpmock.RegisterResponder("GET", "https://lyrist.vercel.app/api/Sonne/Rammstein",
		httpmock.NewStringResponder(http.StatusInternalServerError, ""))

	music, err := dataEnrichmentService.FetchEnrichedMusic(context.Background(), "Rammstein", "Sonne")
	assert.Error(t, err)
	assert.Nil(t, music)

	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info["GET http://ws.audioscrobbler.com/2.0/?method=track.getInfo&api_key=test_api_key&artist=Rammstein&track=Sonne&format=json"])
	assert.Equal(t, 1, info["GET https://lyrist.vercel.app/api/Sonne/Rammstein"])
}

func TestFetchEnrichedMusic_NoLyricsFound(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	cfg := &config.Config{
		APIKey: "test_api_key",
	}

	dataEnrichmentService := service.NewDataEnrichmentService(cfg)

	httpmock.RegisterResponder("GET", "http://ws.audioscrobbler.com/2.0/?method=track.getInfo&api_key=test_api_key&artist=Rammstein&track=Sonne&format=json",
		httpmock.NewStringResponder(http.StatusOK, `{
			"track": {
				"url": "http://www.example.com",
				"wiki": {
					"published": "2 Jan 2006, 15:04"
				}
			}
		}`))

	httpmock.RegisterResponder("GET", "https://lyrist.vercel.app/api/Sonne/Rammstein",
		httpmock.NewStringResponder(http.StatusOK, `{"lyrics": ""}`))

	music, err := dataEnrichmentService.FetchEnrichedMusic(context.Background(), "Rammstein", "Sonne")
	assert.Error(t, err)
	assert.Nil(t, music)
	assert.Equal(t, "no lyrics found for song Sonne by group Rammstein", err.Error())

	info := httpmock.GetCallCountInfo()
	assert.Equal(t, 1, info["GET http://ws.audioscrobbler.com/2.0/?method=track.getInfo&api_key=test_api_key&artist=Rammstein&track=Sonne&format=json"])
	assert.Equal(t, 1, info["GET https://lyrist.vercel.app/api/Sonne/Rammstein"])
}
