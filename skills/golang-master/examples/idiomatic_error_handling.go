// Package examples demonstrates idiomatic Go error handling patterns.
package examples

import (
	"errors"
	"fmt"
	"io"
)

// ─── Sentinel Errors ────────────────────────────────────────────
// Dùng cho known, expected error conditions mà caller cần check.

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrUnauthorized  = errors.New("unauthorized")
)

// ─── Custom Error Types ─────────────────────────────────────────
// Dùng khi cần gắn thêm thông tin vào error.

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field %q %s", e.Field, e.Message)
}

type HTTPError struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

// ─── Error Wrapping ─────────────────────────────────────────────
// Thêm context vào error chain với fmt.Errorf và %w.

func GetUser(id string) (*User, error) {
	user, err := findUserInDB(id)
	if err != nil {
		// Wrap với context mô tả operation đang làm
		return nil, fmt.Errorf("getting user %s: %w", id, err)
	}
	return user, nil
}

func findUserInDB(id string) (*User, error) {
	// Simulate not found
	if id == "" {
		return nil, ErrNotFound
	}
	return &User{ID: id, Name: "Test"}, nil
}

// ─── errors.Is — Check Error Identity ───────────────────────────

func HandleGetUser(id string) {
	user, err := GetUser(id)
	if err != nil {
		// errors.Is kiểm tra trong toàn bộ error chain
		if errors.Is(err, ErrNotFound) {
			fmt.Printf("User %s not found, creating default...\n", id)
			return
		}
		fmt.Printf("Unexpected error: %v\n", err)
		return
	}
	fmt.Printf("Found user: %s\n", user.Name)
}

// ─── errors.As — Extract Error Type ─────────────────────────────

func ProcessRequest() error {
	err := validateRequest("", -1)
	if err != nil {
		var valErr *ValidationError
		if errors.As(err, &valErr) {
			// Truy cập structured fields
			fmt.Printf("Validation failed: field=%s, msg=%s\n", valErr.Field, valErr.Message)
			return err
		}

		var httpErr *HTTPError
		if errors.As(err, &httpErr) {
			fmt.Printf("HTTP error %d: %s\n", httpErr.StatusCode, httpErr.Message)
			return err
		}

		return fmt.Errorf("processing request: %w", err)
	}
	return nil
}

func validateRequest(name string, age int) error {
	if name == "" {
		return &ValidationError{Field: "name", Message: "is required"}
	}
	if age < 0 {
		return &ValidationError{Field: "age", Message: "must be non-negative"}
	}
	return nil
}

// ─── errors.Join — Multiple Errors (Go 1.20+) ──────────────────

func ValidateUser(u User) error {
	var errs []error

	if u.Name == "" {
		errs = append(errs, &ValidationError{Field: "name", Message: "is required"})
	}
	if u.ID == "" {
		errs = append(errs, &ValidationError{Field: "id", Message: "is required"})
	}

	return errors.Join(errs...) // returns nil if errs is empty
}

// ─── Wrapping I/O Errors ────────────────────────────────────────

func ReadConfig(path string) ([]byte, error) {
	f, err := openFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	return data, nil
}

// Helper types
type User struct {
	ID   string
	Name string
}

type file struct{}

func (f *file) Close() error { return nil }

func openFile(path string) (*file, error) {
	if path == "" {
		return nil, ErrNotFound
	}
	return &file{}, nil
}
