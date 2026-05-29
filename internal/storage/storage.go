package storage

import "errors"

type Storage interface {
	SaveIfAbsent(url string, code string) (string, error)
	Get(code string) (string, error)
}

var ErrCodeExists = errors.New("code already exists")
