package service

import (
	"errors"
	"sync"
	"testing"
	"url-shortener/internal/storage"
)

type mockStorage struct {
	mu        sync.RWMutex
	data      map[string]string // code -> url
	urlToCode map[string]string // url -> code
}

func newMockStorage() *mockStorage {
	return &mockStorage{
		data:      make(map[string]string),
		urlToCode: make(map[string]string),
	}
}

func (m *mockStorage) Get(code string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	url, ok := m.data[code]
	if !ok {
		return "", errors.New("not found")
	}
	return url, nil
}

func (m *mockStorage) SaveIfAbsent(url string, code string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if URL already exists
	if existingCode, ok := m.urlToCode[url]; ok {
		return existingCode, nil
	}

	// Check if code already exists
	if _, ok := m.data[code]; ok {
		return "", storage.ErrCodeExists
	}

	// Save both mappings
	m.data[code] = url
	m.urlToCode[url] = code
	return code, nil
}

func TestCreateShortURL_Success(t *testing.T) {
	mock := newMockStorage()
	svc := New(mock)

	code, err := svc.CreateShortURL("https://google.com")
	if err != nil {
		t.Fatalf("CreateShortURL failed: %v", err)
	}
	if len(code) != 10 {
		t.Errorf("Expected code length 10, got %d", len(code))
	}
}

func TestCreateShortURL_ReturnsSameCodeForSameURL(t *testing.T) {
	mock := newMockStorage()
	svc := New(mock)

	code1, err := svc.CreateShortURL("https://google.com")
	if err != nil {
		t.Fatalf("First call failed: %v", err)
	}

	code2, err := svc.CreateShortURL("https://google.com")
	if err != nil {
		t.Fatalf("Second call failed: %v", err)
	}

	if code1 != code2 {
		t.Errorf("Same URL should return same code, got %s and %s", code1, code2)
	}
}

func TestGetOriginalURL_Success(t *testing.T) {
	mock := newMockStorage()
	svc := New(mock)

	code, err := svc.CreateShortURL("https://google.com")
	if err != nil {
		t.Fatalf("CreateShortURL failed: %v", err)
	}

	url, err := svc.GetOriginalURL(code)
	if err != nil {
		t.Fatalf("GetOriginalURL failed: %v", err)
	}

	if url != "https://google.com" {
		t.Errorf("Expected 'https://google.com', got '%s'", url)
	}
}

func TestGetOriginalURL_NotFound(t *testing.T) {
	mock := newMockStorage()
	svc := New(mock)

	_, err := svc.GetOriginalURL("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent code, got nil")
	}
}

func TestConcurrentCreateSameURL(t *testing.T) {
	mock := newMockStorage()
	svc := New(mock)

	iterations := 50
	codes := make([]string, iterations)

	var wg sync.WaitGroup
	wg.Add(iterations)

	for i := 0; i < iterations; i++ {
		go func(idx int) {
			defer wg.Done()
			code, err := svc.CreateShortURL("https://google.com")
			if err != nil {
				t.Errorf("CreateShortURL failed: %v", err)
				return
			}
			codes[idx] = code
		}(i)
	}

	wg.Wait()

	firstCode := codes[0]
	for i := 1; i < iterations; i++ {
		if codes[i] != firstCode {
			t.Errorf("Concurrent requests returned different codes: %s and %s", firstCode, codes[i])
		}
	}
}
