package memory

import (
	"errors"
	"sync"
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

func (m *MemoryStorage) Save(url string, code string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.codeToURL[code] = url
	m.urlToCode[url] = code

	return nil
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

func (m *MemoryStorage) FindByURL(url string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	code, ok := m.urlToCode[url]
	if !ok {
		return "", errors.New("not found")
	}

	return code, nil
}

func (m *MemoryStorage) Exists(code string) (bool, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	_, ok := m.codeToURL[code]

	return ok, nil
}
