package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/basedalex/webcoin/internal/db"
)

type personStore interface {
	GetPerson(ctx context.Context, p db.Person) (db.Person, error)
	CreatePerson(ctx context.Context, p db.Person) (db.Person, error)
}

type Service struct {
	store personStore
}

func New(database personStore) *Service {
	return &Service{
		store: database,
	}
}

func (s *Service) CreatePerson(ctx context.Context, p db.Person) (db.Person, error) {
	person, err := s.store.GetPerson(ctx, p)

	switch {
	case errors.Is(err, db.ErrPersonNotFound):
		person, err = s.store.CreatePerson(ctx, p)
		if err != nil {
			return db.Person{}, fmt.Errorf("error creating person %w", err)
		}
	case err != nil:
		return db.Person{}, fmt.Errorf("error getting person %w", err)
	}

	return person, nil
}
