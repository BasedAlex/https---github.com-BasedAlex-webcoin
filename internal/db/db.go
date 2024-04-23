package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
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
	ID         int       `json:"id"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Patronymic string    `json:"patronymic"`
	Sex        string    `json:"sex"`
	Country    string    `json:"country"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAat"`
}

func (db *Postgres) CreatePerson(ctx context.Context, p Person) (Person, error) {
	stmt := `
	INSERT INTO persons
	(first_name, last_name, patronymic, sex, country, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, first_name, last_name, patronymic, sex, country, created_at, updated_at;`

	rows, err := db.db.Query(ctx, stmt,
		p.FirstName, p.LastName, p.Patronymic,
		p.Sex, p.Country, time.Now(), time.Now())
	if err != nil {
		return Person{}, fmt.Errorf("error in database: %w", err)
	}

	person, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Person])
	if err != nil {
		return Person{}, fmt.Errorf("error collecting rows: %w", err)
	}

	return person, nil
}

func (db *Postgres) FindPerson(ctx context.Context, p Person) (Person, error) {
	stmt := `SELECT 
	(id, first_name, last_name, patronymic, 
	sex, country, created_at, updated_at)
	FROM persons
	WHERE first_name = $1 AND last_name = $2;`

	rows, err := db.db.Query(ctx, stmt, p.FirstName, p.LastName)
	if err != nil {
		person, err := db.CreatePerson(ctx, p)

		return person, err
	}

	person, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Person])
	if err != nil {
		return Person{}, fmt.Errorf("error collecting rows: %w", err)
	}

	return person, nil
}
