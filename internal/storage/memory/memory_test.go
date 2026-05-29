package memory

import (
	"testing"
)

func TestSaveIfAbsent_NewURL_Success(t *testing.T) {
	mem := New()

	code, err := mem.SaveIfAbsent("https://google.com", "abc123")
	if err != nil {
		t.Fatalf("SaveIfAbsent failed: %v", err)
	}
	if code != "abc123" {
		t.Errorf("Expected 'abc123', got '%s'", code)
	}
}

func TestSaveIfAbsent_ExistingURL_ReturnsExistingCode(t *testing.T) {
	mem := New()

	mem.SaveIfAbsent("https://google.com", "abc123")
	code, err := mem.SaveIfAbsent("https://google.com", "xyz789")

	if err != nil {
		t.Fatalf("SaveIfAbsent failed: %v", err)
	}
	if code != "abc123" {
		t.Errorf("Expected existing code 'abc123', got '%s'", code)
	}
}

func TestGet_Success(t *testing.T) {
	mem := New()
	mem.SaveIfAbsent("https://google.com", "abc123")

	url, err := mem.Get("abc123")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if url != "https://google.com" {
		t.Errorf("Expected 'https://google.com', got '%s'", url)
	}
}
