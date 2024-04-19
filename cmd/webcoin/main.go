package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/basedalex/webcoin/internal/config"
	"github.com/basedalex/webcoin/internal/db"
	"github.com/basedalex/webcoin/internal/router"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := config.New()

	dbCtx, cancel := context.WithTimeout(ctx, time.Duration(cfg.Env.DBCancel)*time.Second)
	defer cancel()

	database, err := db.NewPostgres(dbCtx, cfg.Env.DBConn)
	if err != nil {
		log.Println(err)
		log.Exit(1)
	}

	log.Println("connected to db")

	srv := &http.Server{
		Addr:              ":" + cfg.Env.Port,
		Handler:           router.New(database),
		ReadHeaderTimeout: 3 * time.Second,
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*15)

	go func() {
		<-ctx.Done()

		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Print(err)
			os.Exit(1)
		}
	}()

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("error starting server: %s", err)

		return
	}
}
