package examples

import (
	"fmt"
	"strings"
	"testing"
)

// ─── Table-Driven Tests ─────────────────────────────────────────
// Pattern chuẩn trong Go: define test cases as data, loop qua chúng.

func Add(a, b int) int {
	return a + b
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{name: "positive numbers", a: 1, b: 2, expected: 3},
		{name: "zeros", a: 0, b: 0, expected: 0},
		{name: "negative numbers", a: -1, b: -2, expected: -3},
		{name: "mixed signs", a: -1, b: 1, expected: 0},
		{name: "large numbers", a: 1000000, b: 2000000, expected: 3000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

// ─── Subtests với Setup ─────────────────────────────────────────

func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func TestCapitalize(t *testing.T) {
	t.Run("basic cases", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"hello", "Hello"},
			{"world", "World"},
			{"already Capital", "Already Capital"},
		}

		for _, tt := range tests {
			t.Run(tt.input, func(t *testing.T) {
				got := Capitalize(tt.input)
				if got != tt.expected {
					t.Errorf("Capitalize(%q) = %q, want %q", tt.input, got, tt.expected)
				}
			})
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected string
		}{
			{"empty string", "", ""},
			{"single char", "a", "A"},
			{"already uppercase", "A", "A"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := Capitalize(tt.input)
				if got != tt.expected {
					t.Errorf("Capitalize(%q) = %q, want %q", tt.input, got, tt.expected)
				}
			})
		}
	})
}

// ─── Test Helpers ───────────────────────────────────────────────
// Dùng t.Helper() để error messages trỏ đúng call site.

func assertEqual(t *testing.T, got, want interface{}) {
	t.Helper() // Marks this function as a test helper
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestWithHelpers(t *testing.T) {
	result := Add(2, 3)
	assertEqual(t, result, 5) // Nếu fail, error sẽ trỏ đến dòng này, không phải bên trong assertEqual
}

// ─── t.Cleanup — Resource Cleanup ───────────────────────────────

type TestDB struct {
	connected bool
}

func NewTestDB(t *testing.T) *TestDB {
	t.Helper()
	db := &TestDB{connected: true}

	// Cleanup sẽ chạy khi test kết thúc (kể cả khi panic)
	t.Cleanup(func() {
		db.connected = false
		fmt.Println("TestDB cleaned up")
	})

	return db
}

func TestWithCleanup(t *testing.T) {
	db := NewTestDB(t) // cleanup tự động khi test xong
	if !db.connected {
		t.Fatal("expected db to be connected")
	}
}

// ─── Parallel Tests ─────────────────────────────────────────────
// Dùng t.Parallel() cho tests independent có thể chạy đồng thời.

func TestParallel(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"case1", 1, 2, 3},
		{"case2", 10, 20, 30},
		{"case3", -5, 5, 0},
	}

	for _, tt := range tests {
		tt := tt // QUAN TRỌNG: capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Cho phép test này chạy parallel
			got := Add(tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

// ─── Error Case Testing ────────────────────────────────────────

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name      string
		a, b      float64
		expected  float64
		wantError bool
	}{
		{"normal division", 10, 2, 5, false},
		{"division by zero", 10, 0, 0, true},
		{"negative division", -10, 2, -5, false},
		{"fractional result", 1, 3, 0.3333333333333333, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			if tt.wantError {
				assertError(t, err)
				return
			}
			assertNoError(t, err)
			assertEqual(t, got, tt.expected)
		})
	}
}
