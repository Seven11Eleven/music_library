package app

import (
	"fmt"
	"github.com/Seven11Eleven/music_library/api/http/route"
	"github.com/Seven11Eleven/music_library/internal/config"
	"github.com/Seven11Eleven/music_library/internal/database/postgres"
	"github.com/Seven11Eleven/music_library/internal/repository"
	"github.com/Seven11Eleven/music_library/internal/service"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"

	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Router *fiber.App
	DB     *pgxpool.Pool
	Env    *config.Config
}

func NewApp(ctx context.Context) (*App, error) {
	env := config.MustLoad()

	storage, err := postgres.NewDB(*env)
	if storage == nil || err != nil {
		return nil, err
	}

	db := storage.DB()
	if db == nil {
		return nil, fmt.Errorf("failed to establish db conn")
	}

	fiberApp := fiber.New(fiber.Config{
		Immutable: true,
	})

	return &App{
		Router: fiberApp,
		DB:     db,
		Env:    env,
	}, nil
}

func (app *App) Start() {
	musicRepo := repository.NewMusicRepository(app.DB)
	dataEnrichmentService := service.NewDataEnrichmentService(app.Env)
	musicService := service.NewMusicService(musicRepo, dataEnrichmentService)
	route.SetupRoutes(
		app.Router,
		musicService,
		app.Env.ContextTimeout,
	)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("Gracefully Shutting down...")
		app.Router.Shutdown()
	}()
	fmt.Println(app.Env.AppPort)
	if err := app.Router.Listen(":" + app.Env.AppPort); err != nil {
		log.Fatal(err)
	}
}

func (app *App) Close() {
	if app.DB != nil {
		app.DB.Close()
	}

	// Остановка Fiber приложения
	if err := app.Router.Shutdown(); err != nil {
		log.Printf("Error shutting down Fiber: %v", err)
	}

	log.Println("Application closed successfully")
}
