package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/basedalex/webcoin/internal/config"
	"github.com/basedalex/webcoin/internal/db"
	"github.com/basedalex/webcoin/internal/router"
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

	log.Info("connected to db")

	err = router.NewServer(ctx, cfg, database)
	if err != nil {
		log.Panic(err)
	}
}
