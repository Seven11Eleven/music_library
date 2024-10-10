package route

import (
	"github.com/Seven11Eleven/music_library/api/http/controller"
	"github.com/Seven11Eleven/music_library/internal/domain/models"
	"github.com/gofiber/fiber/v2"
	"time"
)

func NewMusicRouter(
	group fiber.Router,
	musicService models.MusicService,
	timeout time.Duration,
) {
	musicController := controller.NewMusicController(musicService)

	group.Get("/info", musicController.GetMusicList)
	group.Get("/verses", musicController.GetVersesOfMusic)
	group.Delete("/:id", musicController.DeleteMusic)
	group.Post("/", musicController.SaveMusic)
	group.Put("/:id", musicController.UpdateMusic)
}
