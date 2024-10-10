package route

import (
	_ "github.com/Seven11Eleven/music_library/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func NewDocsRouter(group fiber.Router) {
	group.Get("/swagger/*", swagger.HandlerDefault)
}
