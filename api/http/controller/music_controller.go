package controller

import (
	"errors"
	"fmt"
	"github.com/Seven11Eleven/music_library/internal/domain/models"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type musicController struct {
	musicService models.MusicService
}

func parseTime(t string) (*time.Time, error) {
	if t == "" {
		log.Warn("Empty music time received")
		return nil, errors.New("music time is empty")
	}

	tm, err := time.Parse("2006-01-02T15:04:05Z", t)
	if err == nil {
		log.Debug("Parsed time successfully using full format")
		return &tm, nil
	}

	tm, err = time.Parse("2006-01-02", t)
	if err != nil {
		log.Errorf("Failed to parse date: %v", err)
		return nil, fmt.Errorf("не удалось распарсить дату: %v", err)
	}

	log.Debug("Parsed time successfully using short format")
	return &tm, nil
}

func NewMusicController(musicService models.MusicService) *musicController {
	log.Info("Creating new music controller instance")
	return &musicController{
		musicService: musicService,
	}
}

func (mc *musicController) GetMusicList(ctx *fiber.Ctx) error {
	log.Info("Fetching music list")
	filters := models.MusicFilters{}

	releaseDateStr := ctx.Query("release_date")
	if releaseDateStr != "" {
		log.Debugf("Received release date filter: %s", releaseDateStr)
		releaseDate, err := parseTime(releaseDateStr)
		if err != nil {
			log.Warnf("Invalid release date format: %v", err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		filters.ReleaseDate = releaseDate
	}

	link := ctx.Query("link")
	if link != "" {
		log.Debugf("Received link filter: %s", link)
		filters.Link = &link
	}

	songName := ctx.Query("song_name")
	if songName != "" {
		log.Debugf("Received song name filter: %s", songName)
		filters.SongName = &songName
	}

	groupName := ctx.Query("group_name")
	if groupName != "" {
		log.Debugf("Received group name filter: %s", groupName)
		filters.GroupName = &groupName
	}

	page := ctx.QueryInt("page")
	pageSize := ctx.QueryInt("page_size")
	log.Debugf("Pagination info: page %d, page_size %d", page, pageSize)

	musicList, err := mc.musicService.GetMusicsByFilters(ctx.Context(), filters, page, pageSize)
	if err != nil {
		log.Errorf("Failed to get music list: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("error: %v", err))
	}

	log.Info("Successfully fetched music list")
	return ctx.JSON(musicList)
}

func (mc *musicController) SaveMusic(ctx *fiber.Ctx) error {
	log.Info("Saving new music")
	req := new(models.MusicQuery)

	if err := ctx.BodyParser(req); err != nil {
		log.Warnf("Failed to parse request body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("error: %v", err))
	}

	savedMusic, err := mc.musicService.SaveMusic(ctx.Context(), req)
	if err != nil {
		log.Errorf("Failed to save music: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("error: %v", err))
	}

	log.Info("Music saved successfully")
	return ctx.JSON(savedMusic)
}

func (mc *musicController) DeleteMusic(ctx *fiber.Ctx) error {
	musicID := ctx.Params("id")
	log.Infof("Deleting music with ID: %s", musicID)
	err := mc.musicService.DeleteMusic(ctx.Context(), musicID)
	if err != nil {
		log.Errorf("Failed to delete music with ID %s: %v", musicID, err)
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("error: %v", err))
	}

	log.Infof("Music with ID %s deleted successfully", musicID)
	return ctx.SendString("Music deleted successfully")
}

func (mc *musicController) GetVersesOfMusic(ctx *fiber.Ctx) error {
	musicID := ctx.Query("music_id")
	log.Infof("Fetching verses for music ID: %s", musicID)

	page := ctx.QueryInt("page")
	pageSize := ctx.QueryInt("page_size")
	log.Debugf("Pagination info: page %d, page_size %d", page, pageSize)

	res, err := mc.musicService.GetMusicTextWithPaginationByVerse(ctx.Context(), musicID, page, pageSize)
	if err != nil {
		log.Errorf("Failed to get verses for music ID %s: %v", musicID, err)
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("error: %v", err))
	}

	log.Infof("Successfully fetched verses for music ID: %s", musicID)
	return ctx.JSON(res)
}

func (mc *musicController) UpdateMusic(ctx *fiber.Ctx) error {
	req := new(models.Music)
	req.ID = ctx.Params("id")
	log.Infof("Updating music with ID: %s", req.ID)

	if err := ctx.BodyParser(req); err != nil {
		log.Warnf("Failed to parse request body: %v", err)
		ctx.Status(http.StatusBadRequest)
		return ctx.SendString(fmt.Sprintf("error: %v", err))
	}

	updatedMusic, err := mc.musicService.UpdateMusic(ctx.Context(), *req)
	if err != nil {
		log.Errorf("Failed to update music with ID %s: %v", req.ID, err)
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.SendString(fmt.Sprintf("error: %v", err))
	}

	log.Infof("Music with ID %s updated successfully", req.ID)
	return ctx.JSON(updatedMusic)
}
