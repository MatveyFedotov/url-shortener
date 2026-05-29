package service

import (
	"errors"
	"url-shortener/internal/generator"
	"url-shortener/internal/storage"
)

type Storage interface {
	Get(code string) (string, error)
	SaveIfAbsent(url string, code string) (string, error)
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
	for {
		code := generator.GenerateCode(10)

		savedCode, err := s.storage.SaveIfAbsent(originalURL, code)
		if errors.Is(err, storage.ErrCodeExists) {
			continue
		}
		if err != nil {
			return "", err
		}

		return savedCode, nil
	}
}

func (s *Service) GetOriginalURL(code string) (string, error) {
	return s.storage.Get(code)
}
