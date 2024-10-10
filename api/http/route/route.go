package route

import (
	"github.com/Seven11Eleven/music_library/api/http/middleware"
	"github.com/Seven11Eleven/music_library/internal/domain/models"
	"github.com/gofiber/fiber/v2"
	"time"
)

func SetupRoutes(
	app *fiber.App,
	musicService models.MusicService,
	timeout time.Duration,
) {
	middleware.MiddlewaresSetup(app)

	musicRoute := app.Group("/music")
	NewMusicRouter(musicRoute, musicService, timeout)

	docsRoute := app.Group("/docs")
	NewDocsRouter(docsRoute)
}
