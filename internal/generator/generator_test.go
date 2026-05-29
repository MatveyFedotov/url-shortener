package generator

import (
	"testing"
)

func TestGenerateCode_Length(t *testing.T) {
	code := GenerateCode(10)
	if len(code) != 10 {
		t.Errorf("Expected length 10, got %d", len(code))
	}
}

func TestGenerateCode_OnlyAllowedChars(t *testing.T) {
	code := GenerateCode(1000)
	allowedChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

	for _, ch := range code {
		found := false
		for _, allowed := range allowedChars {
			if ch == allowed {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Invalid character '%c' found in generated code", ch)
		}
	}
}

func TestGenerateCode_Uniqueness(t *testing.T) {
	codes := make(map[string]bool)
	iterations := 10000

	for i := 0; i < iterations; i++ {
		code := GenerateCode(10)
		if codes[code] {
			t.Errorf("Duplicate code generated: %s", code)
		}
		codes[code] = true
	}
}

func TestGenerateCode_Randomness(t *testing.T) {
	code1 := GenerateCode(50)
	code2 := GenerateCode(50)

	if code1 == code2 {
		t.Error("Two consecutive calls returned the same string, expected randomness")
	}
}

func BenchmarkGenerateCode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateCode(10)
	}
}
