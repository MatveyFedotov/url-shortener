package service

import (
	"testing"
	"url-shortener/internal/storage/memory"
)

func TestCreateShortURL_ValidURL(t *testing.T) {
	store := memory.New()
	svc := New(store)

	code, err := svc.CreateShortURL("https://google.com")
	if err != nil {
		t.Fatalf("Ошибка CreateShortURL: %v", err)
	}

	if len(code) != 10 {
		t.Errorf("Ожидалась длина кода 10, получено %d", len(code))
	}
}

func TestCreateShortURL_InvalidURL(t *testing.T) {
	store := memory.New()
	svc := New(store)

	_, err := svc.CreateShortURL("not-a-url")
	if err == nil {
		t.Error("Ожидалась ошибка для невалидного URL")
	}
}

func TestCreateShortURL_Uniqueness(t *testing.T) {
	store := memory.New()
	svc := New(store)

	code1, _ := svc.CreateShortURL("https://google.com")
	code2, _ := svc.CreateShortURL("https://google.com")

	if code1 != code2 {
		t.Errorf("Один URL должен возвращать один код. Получены %s и %s", code1, code2)
	}
}

func TestGetOriginalURL_Success(t *testing.T) {
	store := memory.New()
	svc := New(store)

	code, _ := svc.CreateShortURL("https://google.com")
	url, err := svc.GetOriginalURL(code)

	if err != nil {
		t.Fatalf("Ошибка GetOriginalURL: %v", err)
	}

	if url != "https://google.com" {
		t.Errorf("Ожидался 'https://google.com', получено '%s'", url)
	}
}

func TestGetOriginalURL_NotFound(t *testing.T) {
	store := memory.New()
	svc := New(store)

	_, err := svc.GetOriginalURL("notexist")
	if err == nil {
		t.Error("Ожидалась ошибка для несуществующего кода")
	}
}
