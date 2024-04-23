package service

import (
	"context"

	"github.com/basedalex/webcoin/internal/db"
)

type personStore interface {
	CreatePerson(ctx context.Context, p db.Person) (db.Person, error)
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

	newperson, err := s.Database.CreatePerson(ctx, p)

	return newperson, err
}

