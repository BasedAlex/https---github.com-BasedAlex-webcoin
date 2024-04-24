package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/basedalex/webcoin/internal/config"
	"github.com/basedalex/webcoin/internal/db"
	"github.com/basedalex/webcoin/internal/router"
	"github.com/basedalex/webcoin/internal/service"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := config.New()

	database, err := db.NewPostgres(ctx, cfg.Env.PGDSN)
	if err != nil {
		log.Panic(err)
	}

	service := service.New(database)

	log.Info("connected to db")

	err = router.NewServer(ctx, cfg, service)
	if err != nil {
		log.Panic(err)
	}
}
