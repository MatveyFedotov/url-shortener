package generator

import (
	"testing"
)

func TestGenerateCode_Length(t *testing.T) {
	code := GenerateCode(10)
	if len(code) != 10 {
		t.Errorf("Ожидалась длина 10, получено %d", len(code))
	}
}

func TestGenerateCode_AllowedChars(t *testing.T) {
	code := GenerateCode(100)
	allowed := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

	for _, ch := range code {
		found := false
		for _, a := range allowed {
			if ch == a {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Символ '%c' не разрешён", ch)
		}
	}
}

func TestGenerateCode_Uniqueness(t *testing.T) {
	codes := make(map[string]bool)
	for i := 0; i < 100; i++ {
		code := GenerateCode(10)
		if codes[code] {
			t.Errorf("Найден дубликат: %s", code)
		}
		codes[code] = true
	}
}
