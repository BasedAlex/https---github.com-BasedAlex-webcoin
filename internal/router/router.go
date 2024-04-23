package router

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
		}
	}()

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("error with the server: %w", err)
	}

	return nil
}

type personStore interface {
	CreatePerson(ctx context.Context, p db.Person) (db.Person, error)
}

type Handler struct {
	Database personStore
}

type HTTPResponse struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func newRouter(database personStore) *chi.Mux {
	handler := &Handler{
		Database: database,
	}

	r := chi.NewRouter()

	r.Post("/api/v1/person", handler.createPerson)

	return r
}

func (h *Handler) createPerson(w http.ResponseWriter, r *http.Request) {
	var person db.Person

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		writeErrResponse(w, err, http.StatusBadRequest)

		return
	}

	newperson, err := h.Database.CreatePerson(r.Context(), person)
	if err != nil {
		writeErrResponse(w, err, http.StatusInternalServerError)

		return
	}

	writeOkResponse(w, http.StatusCreated, newperson)
}

func writeOkResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(HTTPResponse{Data: data})
	if err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func writeErrResponse(w http.ResponseWriter, err error, statusCode int) {
	log.Error(err)
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, err.Error(), statusCode)
}
