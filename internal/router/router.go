package router

import (
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
	
}