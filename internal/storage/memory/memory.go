package memory

import (
	"errors"
	"sync"
	"url-shortener/internal/storage"
)

type MemoryStorage struct {
	codeToURL map[string]string
	urlToCode map[string]string
	mutex     sync.RWMutex
}

func New() *MemoryStorage {
	return &MemoryStorage{
		codeToURL: make(map[string]string),
		urlToCode: make(map[string]string),
	}
}

func (m *MemoryStorage) SaveIfAbsent(url string, code string) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if existingCode, ok := m.urlToCode[url]; ok {
		return existingCode, nil
	}

	if _, ok := m.codeToURL[code]; ok {
		return "", storage.ErrCodeExists
	}
	m.codeToURL[code] = url
	m.urlToCode[url] = code

	return code, nil
}

func (m *MemoryStorage) Get(code string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	url, ok := m.codeToURL[code]
	if !ok {
		return "", errors.New("not found")
	}

	return url, nil
}
