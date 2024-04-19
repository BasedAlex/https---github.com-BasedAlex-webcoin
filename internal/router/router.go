package router

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/basedalex/webcoin/internal/db"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Database *db.Postgres
}

func New(database *db.Postgres) *chi.Mux {
	handler := &Handler{
		Database: database,
	}
	r := chi.NewRouter()

	r.Post("/person", handler.postPerson)

	return r
}

func (h *Handler) postPerson(w http.ResponseWriter, r *http.Request) {
	var person db.Person

	err := json.NewDecoder(r.Body).Decode(&person)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	ctx := context.TODO()

	err = h.Database.CreatePerson(ctx, person)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	w.WriteHeader(http.StatusCreated)
}