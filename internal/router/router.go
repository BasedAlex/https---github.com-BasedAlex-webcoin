package router

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/basedalex/webcoin/internal/config"
	"github.com/basedalex/webcoin/internal/db"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

func NewServer(ctx context.Context, cfg *config.Config, database *db.Postgres) error {
	srv := &http.Server{
		Addr:              ":" + cfg.Env.Port,
		Handler:           newRouter(database),
		ReadHeaderTimeout: 3 * time.Second,
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*15)

	go func() {
		<-ctx.Done()

		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Warn(err)
			os.Exit(1)
		}
	}()

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error with the server: %w", err)
	}

	return nil
}

type personStore interface {
	CreatePerson(ctx context.Context, p db.Person) error
}

type Handler struct {
	Database personStore
}

func newRouter(database *db.Postgres) *chi.Mux {
	handler := &Handler{
		Database: database,
	}

	r := chi.NewRouter()

	r.Get("/ping", ping)

	r.Post("/person", handler.createPerson)

	return r
}

func ping(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "hello")
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) createPerson(w http.ResponseWriter, r *http.Request) {
	var person db.Person

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = h.Database.CreatePerson(r.Context(), person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
}
