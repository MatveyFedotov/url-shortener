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

func (p *PostgresStorage) SaveIfAbsent(url string, code string) (string, error) {
	var savedCode string

	err := p.db.QueryRow(context.Background(), `
  INSERT INTO links (url, code)
  VALUES ($1, $2)
  ON CONFLICT (url) DO UPDATE SET url = links.url
  RETURNING code
 `, url, code).Scan(&savedCode)

	if err != nil {
		return "", err
	}

	return savedCode, nil
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
