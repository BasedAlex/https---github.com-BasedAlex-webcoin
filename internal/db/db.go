package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Postgres struct {
	db *pgx.Conn
}

func NewPostgres(ctx context.Context, dbConnect string) (*Postgres, error) {
	db, err := pgx.Connect(ctx, dbConnect)
	if err != nil {
		return nil, err
	}
	return &Postgres{db: db}, nil
}


type Person struct {
	Id int `json:"int"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Patronymic string `json:"patronymic"`
	Sex string `json:"sex"`
	Country string `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (db *Postgres) CreatePerson(ctx context.Context, person Person) error {
	
	stmt := `
	INSERT INTO persons
	(first_name, last_name, patronymic, sex, country, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7);`

	_, err := db.db.Exec(ctx, stmt, person.FirstName, person.LastName, person.Patronymic, person.Sex, person.Country, time.Now(), time.Now())

	if err != nil {
		return err
	}

	return nil
}