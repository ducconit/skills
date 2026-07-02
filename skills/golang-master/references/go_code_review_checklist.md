# Go Code Review Checklist

Checklist dùng khi review code Go. Đảm bảo tất cả các mục đều pass trước khi approve.

## Error Handling

- [ ] Tất cả errors đều được kiểm tra (`if err != nil`)
- [ ] Errors được wrap với context (`fmt.Errorf("doing x: %w", err)`)
- [ ] Error messages viết lowercase, không có dấu chấm cuối
- [ ] Sentinel errors dùng cho known conditions
- [ ] Custom error types khi cần thêm thông tin

## Concurrency

- [ ] Không có goroutine leaks (luôn có cách exit)
- [ ] Proper context cancellation và propagation
- [ ] Channels được close ở producer side
- [ ] sync.Mutex dùng đúng chỗ, defer Unlock()
- [ ] Chạy tests với `-race` flag

## Design

- [ ] Interfaces nhỏ (1-3 methods), defined ở consumer side
- [ ] Không có premature abstraction
- [ ] Accept interfaces, return structs
- [ ] Zero value của struct có ý nghĩa
- [ ] Package dependencies không circular

## Code Quality

- [ ] Comments giải thích WHY, không phải WHAT
- [ ] Exported functions/types có doc comments
- [ ] Không có magic numbers (dùng constants)
- [ ] Naming theo Go conventions (camelCase, acronyms ALL_CAPS)
- [ ] Proper use of `defer` (resource cleanup, mutex unlock)

## Testing

- [ ] Tests cover happy path VÀ edge cases
- [ ] Table-driven tests cho hàm có nhiều cases
- [ ] Test helpers dùng `t.Helper()`
- [ ] Cleanup resources với `t.Cleanup()`
- [ ] Không có flaky tests (no timing dependencies)

## Performance

- [ ] Slices preallocated khi biết capacity
- [ ] strings.Builder cho string concatenation
- [ ] Không allocation không cần thiết trong hot paths
- [ ] Proper use of pointers vs values
- [ ] bufio cho I/O operations

## Security

- [ ] Input validation cho tất cả user input
- [ ] Parameterized queries cho database
- [ ] Không hardcode secrets/credentials
- [ ] TLS cho network communications
- [ ] Không log sensitive data
