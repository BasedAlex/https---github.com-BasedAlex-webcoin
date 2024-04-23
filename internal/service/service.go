package service

import (
	"context"
	"fmt"

	"github.com/basedalex/webcoin/internal/db"
)

type personStore interface {
	FindPerson(ctx context.Context, p db.Person) (db.Person, error)
	// CreatePerson(ctx context.Context, p db.Person) (db.Person, error)
}

type Service struct {
	Database personStore
}

func New(database personStore) *Service {
	return &Service{
		Database: database,
	}
}

func (s *Service) CreatePerson(ctx context.Context, p db.Person) (db.Person, error) {
	newperson, err := s.Database.FindPerson(ctx, p)

	return newperson, fmt.Errorf("error creating person %w", err)
}
