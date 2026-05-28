package service

import (
	"errors"
	"net/url"

	"url-shortener/internal/generator"
)

type Storage interface {
	Save(url string, code string) error
	Get(code string) (string, error)
	FindByURL(url string) (string, error)
	Exists(code string) (bool, error)
}

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) CreateShortURL(originalURL string) (string, error) {
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return "", errors.New("invalid url")
	}

	existingCode, err := s.storage.FindByURL(originalURL)
	if err == nil {
		return existingCode, nil
	}

	var code string

	for {
		code = generator.GenerateCode(10)

		exists, err := s.storage.Exists(code)
		if err != nil {
			return "", err
		}

		if !exists {
			break
		}
	}

	err = s.storage.Save(originalURL, code)
	if err != nil {
		return "", err
	}

	return code, nil
}

func (s *Service) GetOriginalURL(code string) (string, error) {
	return s.storage.Get(code)
}
