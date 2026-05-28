package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	db *pgxpool.Pool
}

func New(conn string) (*PostgresStorage, error) {
	db, err := pgxpool.New(context.Background(), conn)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

func (p *PostgresStorage) Save(url string, code string) error {
	_, err := p.db.Exec(context.Background(),
		"INSERT INTO links (url, code) VALUES ($1, $2)",
		url, code,
	)
	return err
}

func (p *PostgresStorage) Get(code string) (string, error) {
	var url string

	err := p.db.QueryRow(context.Background(),
		"SELECT url FROM links WHERE code=$1",
		code,
	).Scan(&url)

	if err != nil {
		return "", errors.New("not found")
	}

	return url, nil
}

func (p *PostgresStorage) FindByURL(url string) (string, error) {
	var code string

	err := p.db.QueryRow(context.Background(),
		"SELECT code FROM links WHERE url=$1",
		url,
	).Scan(&code)

	if err != nil {
		return "", errors.New("not found")
	}

	return code, nil
}

func (p *PostgresStorage) Exists(code string) (bool, error) {
	var exists bool

	err := p.db.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM links WHERE code=$1)",
		code,
	).Scan(&exists)

	return exists, err
}
