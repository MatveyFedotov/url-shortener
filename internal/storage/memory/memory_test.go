package memory

import (
	"testing"
)

func TestMemoryStorage_SaveAndGet(t *testing.T) {
	store := New()

	err := store.Save("https://google.com", "abc123")
	if err != nil {
		t.Fatalf("Ошибка Save: %v", err)
	}

	url, err := store.Get("abc123")
	if err != nil {
		t.Fatalf("Ошибка Get: %v", err)
	}

	if url != "https://google.com" {
		t.Errorf("Ожидался 'https://google.com', получено '%s'", url)
	}
}

func TestMemoryStorage_GetNotFound(t *testing.T) {
	store := New()

	_, err := store.Get("notexist")
	if err == nil {
		t.Error("Ожидалась ошибка для несуществующего кода")
	}
}

func TestMemoryStorage_FindByURL(t *testing.T) {
	store := New()
	store.Save("https://google.com", "abc123")

	code, err := store.FindByURL("https://google.com")
	if err != nil {
		t.Fatalf("Ошибка FindByURL: %v", err)
	}

	if code != "abc123" {
		t.Errorf("Ожидался 'abc123', получено '%s'", code)
	}
}

func TestMemoryStorage_Exists(t *testing.T) {
	store := New()
	store.Save("https://google.com", "abc123")

	exists, _ := store.Exists("abc123")
	if !exists {
		t.Error("Ожидалось exists=true для существующего кода")
	}

	exists, _ = store.Exists("notexist")
	if exists {
		t.Error("Ожидалось exists=false для несуществующего кода")
	}
}
