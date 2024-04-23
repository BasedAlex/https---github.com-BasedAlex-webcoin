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

func NewServer(ctx context.Context, cfg *config.Config, service serviceLayer) error {
	srv := &http.Server{
		Addr:              ":" + cfg.Env.Port,
		Handler:           newRouter(service),
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

type serviceLayer interface {
	// FindPerson(ctx, p db.Person) (db.Person, error)
	CreatePerson(ctx context.Context, p db.Person) (db.Person, error)
}

type Handler struct {
	Service serviceLayer
}

type HTTPResponse struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func newRouter(service serviceLayer) *chi.Mux {
	handler := &Handler{
		Service: service,
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

	newperson, err := h.Service.CreatePerson(r.Context(), person)
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	log.Warn(err)

	jsonErr := json.NewEncoder(w).Encode(HTTPResponse{Error: err.Error()})
	if jsonErr != nil {
		log.Warn(jsonErr)
	}
}
