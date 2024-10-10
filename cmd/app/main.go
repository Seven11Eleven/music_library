package main

import (
	"context"
	"github.com/Seven11Eleven/music_library/internal/app"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	context, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	newApp, err := app.NewApp(context)
	if err != nil {
		log.Fatal(err)
	}
	defer newApp.Close()

	newApp.Start()

}
