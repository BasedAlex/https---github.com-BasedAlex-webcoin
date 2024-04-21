package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	db *pgxpool.Pool
}

func NewPostgres(ctx context.Context, dbConnect string) (*Postgres, error) {
	config, err := pgxpool.ParseConfig(dbConnect)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %w", err)
	}

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("error pinging to database: %w", err)
	}

	return &Postgres{db: db}, nil
}

type Person struct {
	ID         int       `json:"int"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Patronymic string    `json:"patronymic"`
	Sex        string    `json:"sex"`
	Country    string    `json:"country"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAat"`
}

func (db *Postgres) CreatePerson(ctx context.Context, p Person) error {
	stmt := `
	INSERT INTO persons
	(first_name, last_name, patronymic, sex, country, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7);`

	_, err := db.db.Exec(ctx, stmt, p.FirstName, p.LastName, p.Patronymic, p.Sex, p.Country, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("error in database: %w", err)
	}

	return nil
}
