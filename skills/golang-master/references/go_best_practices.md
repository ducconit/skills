# Go Best Practices — Tham chiếu Toàn diện

## 1. Code Organization & Project Structure

Lựa chọn mô hình tổ chức mã nguồn có ảnh hưởng trực tiếp đến khả năng bảo trì và kiểm soát lỗi phụ thuộc vòng (**Circular Dependency**). Có hai mô hình chính thường dùng:

### Mô hình 1: Standard Go Project Layout
Tập trung vào ranh giới kiến trúc nghiêm ngặt bằng các thư mục đặc biệt của Go. Phù hợp cho các dự án lớn, doanh nghiệp.
- **`cmd/`**: Điểm khởi chạy (wiring), khởi tạo hệ thống và tiêm phụ thuộc (Dependency Injection), tuyệt đối không chứa logic nghiệp vụ.
- **`internal/`**: Ranh giới bảo mật của Go. Mã đặt tại đây chỉ module sở hữu mới có thể import. Thường tổ chức theo mô hình lai 5 tầng:
  - `internal/handler/`: Bộ chuyển đổi HTTP/gRPC siêu mỏng, parse request/response.
  - `internal/service/`: Lớp trung tâm thực thi nghiệp vụ (Business logic).
  - `internal/repository/`: Chứa các truy vấn SQL vật lý, giao tiếp DB.
  - `internal/domain/`: Định nghĩa struct thực thể, core interfaces, và mã lỗi (không import package nào khác).
- **`pkg/`**: Mã tiện ích dùng chung có thể import tự do từ các module ngoài.

### Mô hình 2: Standard Package Layout (Ben Johnson - WTF Dial)
Triết lý tối giản, tập trung vào **Domain** thay vì tầng kỹ thuật. Cấu trúc thư mục được thiết kế phẳng:
- **Gói gốc (Root Package - `/wtf`)**: Đại diện cho Domain nghiệp vụ. Định nghĩa toàn bộ Domain Models (struct như `User`, `Dial`) và Interfaces trừu tượng (`UserService`, `DialService`).
- **Các gói con (Subpackages)**: Đóng vai trò là các Adapters triển khai kỹ thuật cụ thể (ví dụ: `/postgres` chứa triển khai PostgreSQL cho interface ở root, `/http` chứa HTTP Handlers).
- Các gói con phụ thuộc 1 chiều vào root package, nghiêm cấm các gói con import lẫn nhau (giao tiếp chéo phải qua interfaces ở root).

### Cách giải quyết lỗi Circular Dependency trong Go
1. **Gom nhóm theo Tính năng (Feature-based)** thay vì tầng kỹ thuật (Layer-based): Gom tất cả tệp liên quan đến một miền (handler, service, db của `user`) vào chung package `user` để gọi nhau tự do mà không cần import chéo.
2. **Vùng lõi không phụ thuộc**: Thiết lập package `internal/domain` (ở mô hình 1) hoặc root package (ở mô hình 2) là nơi chứa structs và interfaces lõi không import bất kỳ package nội bộ nào khác. Tất cả các package khác import nó 1 chiều.
3. **Định nghĩa Interface ở nơi sử dụng (Client-side)**, không định nghĩa ở nơi triển khai (Provider-side). Giúp bẻ gãy quan hệ phụ thuộc chéo nhờ cơ chế Implicit Interfaces của Go.
4. **Tập trung hóa wiring tại `main.go`**: Đưa toàn bộ code khởi tạo và tiêm phụ thuộc thủ công hoặc qua công cụ (Google Wire, Uber Fx) về file `cmd/.../main.go`.

## 2. Naming Conventions

```go
// Package names: lowercase, singular, no underscores
package user     // ✅
package users    // ❌
package user_mgr // ❌

// Variables: camelCase, ngắn gọn, ý nghĩa
var userCount int           // ✅
var numberOfUsers int       // ❌ quá dài cho local var
var n int                   // ✅ cho loop/short-lived

// Interfaces: verb + "er" suffix cho single-method
type Reader interface { Read(p []byte) (n int, err error) }
type Validator interface { Validate() error }

// Exported: PascalCase, unexported: camelCase
func ProcessOrder() {}      // ✅ exported
func processOrder() {}      // ✅ unexported

// Acronyms: all caps
var httpClient *http.Client // ✅
var userID string           // ✅
var userId string           // ❌

// Getters: không dùng "Get" prefix
func (u *User) Name() string {}    // ✅
func (u *User) GetName() string {} // ❌
```

## 3. Error Handling

### Nguyên tắc cốt lõi

```go
// Luôn kiểm tra error
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doing something: %w", err) // wrap với context
}

// Sentinel errors cho known conditions
var ErrNotFound = errors.New("not found")
var ErrAlreadyExists = errors.New("already exists")

// Custom error types cho rich error info
type ValidationError struct {
    Field   string
    Message string
}
func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Message)
}

// Kiểm tra error type
if errors.Is(err, ErrNotFound) { /* handle */ }
var valErr *ValidationError
if errors.As(err, &valErr) { /* handle with valErr.Field */ }

// Gộp nhiều lỗi bằng errors.Join (Go 1.20+) thay vì tự viết logic
var errs error
if err := validateUsername(u); err != nil {
    errs = errors.Join(errs, err)
}
if err := validateEmail(e); err != nil {
    errs = errors.Join(errs, err)
}

// Error messages: lowercase, no punctuation
return fmt.Errorf("opening file: %w", err)        // ✅
return fmt.Errorf("Failed to open file: %w", err)  // ❌
```

### Patterns nên tránh

```go
// ❌ Bỏ qua error
result, _ := doSomething()

// ❌ Chỉ log mà không return/handle
if err != nil {
    log.Println(err)
}

// ❌ Wrap error thiếu context
return fmt.Errorf("%w", err) // không có context

// ❌ Double wrapping
return fmt.Errorf("failed: %w", fmt.Errorf("error: %w", err))
```

### Anti-pattern: Truyền dependencies qua `context.Value`
`context.Value` chỉ dành cho **request-scoped data** (request ID, trace ID, auth token). Tuyệt đối không dùng để truyền dependencies (DB, config, logger) vì nó ẩn coupling, khó test, khó debug.

```go
// ❌ Anti-pattern: Dependencies ẩn trong context
ctx = context.WithValue(ctx, "db", db)
db := ctx.Value("db").(*sql.DB) // type assertion nguy hiểm, nil panic tiềm ẩn

// ✅ Explicit dependency injection qua constructor
func NewHandler(db *sql.DB) *Handler {
    return &Handler{db: db}
}
```

## 4. Concurrency Patterns

### Goroutines & Channels

```go
// Luôn có cơ chế dừng goroutine
func worker(ctx context.Context, jobs <-chan Job) {
    for {
        select {
        case <-ctx.Done():
            return
        case job, ok := <-jobs:
            if !ok {
                return
            }
            process(job)
        }
    }
}

// Đóng channel ở producer side, không phải consumer
func produce(ch chan<- int) {
    defer close(ch)
    for i := 0; i < 10; i++ {
        ch <- i
    }
}
```

### sync Package

```go
// sync.Mutex cho critical sections
type SafeCounter struct {
    mu sync.Mutex
    v  map[string]int
}
func (c *SafeCounter) Inc(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.v[key]++
}

// sync.RWMutex khi read nhiều hơn write
// sync.Once cho one-time initialization
// sync.WaitGroup cho waiting goroutines
// sync.Pool cho object reuse
```

### errgroup cho concurrent operations
Dùng `golang.org/x/sync/errgroup` để quản lý nhiều goroutine chạy đồng thời, tự động cancel context khi một goroutine gặp lỗi.

```go
// errgroup tự động lan truyền việc hủy bỏ (cancel) context khi có lỗi
g, ctx := errgroup.WithContext(ctx)
for _, url := range urls {
    url := url // Tránh biến lặp bị ghi đè (Go 1.22+ đã tự động fix, nhưng vẫn nên chú ý ở các phiên bản cũ)
    g.Go(func() error {
        return fetch(ctx, url)
    })
}
if err := g.Wait(); err != nil {
    return err
}
```

### Structured Concurrency & Panic Propagation (Thư viện `sourcegraph/conc`)
Trong Go, nếu goroutine con bị panic, toàn bộ ứng dụng sẽ bị sập. Dùng `conc.WaitGroup` để tự động catch panic ở goroutine con và re-panic về luồng chính:

```go
import "github.com/sourcegraph/conc"

var wg conc.WaitGroup
for _, job := range jobs {
    job := job
    wg.Go(func() {
        process(job) // Nếu process() panic, nó sẽ được bắt lại ở Wait() và re-panic ở luồng chính
    })
}
wg.Wait() // An toàn, tránh panic sập toàn bộ app thầm lặng
```

### Quy trình Graceful Shutdown chuẩn
1. Lắng nghe tín hiệu hệ thống (`SIGINT`, `SIGTERM`) qua signal channel.
2. Khi nhận tín hiệu, broadcast cancel context để báo cho các workers dừng xử lý công việc mới.
3. Gọi `server.Shutdown(ctx)` với thời gian timeout hợp lý (ví dụ: 10s-30s) để đợi các request hiện tại hoàn thành.
4. Đóng các tài nguyên nền (DB Connection, RabbitMQ, Logger flush).
5. Tuyệt đối tránh gọi `os.Exit()` tùy tiện vì nó sẽ bỏ qua mọi khối lệnh `defer`.

## 5. Interface Design

```go
// Interfaces nhỏ, tập trung
type Writer interface {
    Write(p []byte) (n int, err error)
}

// Accept interfaces, return structs
func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

// Define interface ở consumer side, không phải producer
// package handler
type UserStore interface {
    GetUser(ctx context.Context, id string) (*User, error)
}
type Handler struct { store UserStore }

// Không tạo interface "trước khi cần"
// Interface nên phát sinh từ nhu cầu thực tế (testing, multiple implementations)

### Gotcha: Kế thừa ảo (Virtual Inheritance) trong Struct Embedding
Struct embedding giúp tái sử dụng code thông qua composition. Tuy nhiên, Go KHÔNG hỗ trợ virtual inheritance. Phương thức của struct cha gọi một phương thức khác sẽ KHÔNG gọi phương thức bị override ở struct con:

```go
type Parent struct{}
func (p *Parent) Set() { p.Do() }
func (p *Parent) Do()  { fmt.Println("Parent Do") }

type Child struct { Parent }
func (c *Child) Do() { fmt.Println("Child Do") }

// Thực thi
c := Child{}
c.Set() // Output: "Parent Do" (Không phải "Child Do" do c.Parent.Do() được gọi trực tiếp)
```
```

## 6. Testing

```go
// Table-driven tests
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 1, 2, 3},
        {"zero", 0, 0, 0},
        {"negative", -1, 1, 0},
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

// Test helpers
func setupTestDB(t *testing.T) *DB {
    t.Helper()
    db := NewTestDB()
    t.Cleanup(func() { db.Close() })
    return db
}

// Parallel tests
t.Run("parallel_test", func(t *testing.T) {
    t.Parallel()
    // test code
})
```

## 7. Performance

- **Đo trước khi tối ưu**: Dùng `go test -bench` và `pprof`
- **Preallocate slices**: `make([]T, 0, expectedCap)`
- **strings.Builder** cho string concatenation
- **sync.Pool** cho objects tạo/hủy thường xuyên
- **bufio** cho I/O operations
- **Tránh reflection** trong hot paths
- **Tránh allocations** không cần thiết (pointer vs value receiver)
- Xem chi tiết tại `references/go_performance_patterns.md`

## 8. API Design

### Functional Options

```go
type Option func(*Server)

func WithPort(port int) Option {
    return func(s *Server) { s.port = port }
}

func NewServer(opts ...Option) *Server {
    s := &Server{port: 8080} // defaults
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
srv := NewServer(WithPort(9090), WithTimeout(30*time.Second))
```

## 9. Dependency Management

- Luôn commit `go.sum`
- Dùng `go mod tidy` thường xuyên
- Pin major versions
- Kiểm tra vulnerabilities: `go install golang.org/x/vuln/cmd/govulncheck@latest && govulncheck ./...`
- Tránh quá nhiều dependencies

## 10. Logging & Observability

Từ Go 1.21+, thư viện chuẩn bổ sung package `log/slog` hỗ trợ ghi log có cấu trúc hiệu năng cao, giảm phụ thuộc vào các thư viện bên thứ ba (như `zerolog` hay `zap`).

### Tùy chọn 1: log/slog (Standard Library - Khuyến nghị cho dự án hiện đại)
```go
import "log/slog"

// JSON Structured logging
slog.Info("user created",
    slog.String("user_id", user.ID),
    slog.Int("age", user.Age),
)

// Context-aware/Grouped logging
logger := slog.Default().With(slog.String("request_id", requestID))
logger.Info("processing request")

// Error logging
slog.Error("failed to create user",
    slog.Any("err", err),
    slog.String("operation", "create_user"),
)
```

### Tùy chọn 2: zerolog (High Performance)
Dành cho các hệ thống cực kỳ nhạy cảm về hiệu năng và phân bổ bộ nhớ (Zero allocation logger) — xem chi tiết tại `references/go_recommended_stack.md`.
```go
import (
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

log.Info().
    Str("user_id", user.ID).
    Int("age", user.Age).
    Msg("user created")
```

## 11. Modern HTTP Routing (Go 1.22+)

Go 1.22+ cải tiến toàn diện bộ định tuyến gốc `net/http.ServeMux` hỗ trợ HTTP Verbs và Path Parameters:

```go
mux := http.NewServeMux()
// Khớp phương thức GET và lấy path parameter {id}
mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintf(w, "User ID: %s", id)
})

// Khớp chính xác phần kết thúc đường dẫn bằng {$
mux.HandleFunc("GET /users/{$}", listUsersHandler)
```

### Cạm bẫy (Gotchas) của ServeMux mới:
1. **Không hỗ trợ biểu thức chính quy (Regex)** trên biến đường dẫn (như `go-chi/chi` hay `gin`). Mọi việc xác thực kiểu dữ liệu của parameter (ví dụ: bắt buộc `{id}` phải là số) đều phải thực hiện thủ công bằng code trong handler.
2. **Không có Middleware Grouping**: ServeMux không hỗ trợ gán middleware theo nhóm router con (như `.Group()` hay `.Route()`). Cần viết custom adapter hoặc kết hợp với router siêu nhẹ như `go-chi/chi` nếu hệ thống phức tạp.

## 12. Iterators (Go 1.23+)

Go 1.23+ giới thiệu cơ chế lặp chuẩn hóa (Range-over-func) với hai signature chính:
- `type Seq[V any] func(yield func(V) bool)` (lặp 1 giá trị)
- `type Seq2[K, V any] func(yield func(K, V) bool)` (lặp khóa/giá trị)

### Gotcha trí mạng: Bắt buộc kiểm tra kết quả trả về của `yield`
Hàm `yield` trả về `false` khi vòng lặp của client bị dừng sớm (gặp `break`, `return` hoặc `panic`). Iterator **phải dừng ngay lập tức** khi phát hiện `yield` trả về `false` để tránh rò rỉ tài nguyên.

```go
// Triển khai lặp đệ quy an toàn cho Cây Nhị Phân (ngăn mạch ngay khi yield = false)
func (t *Tree[E]) push(yield func(E) bool) bool {
    if t == nil {
        return true
    }
    // Sử dụng toán tử && để ngắt lặp đệ quy ngay lập tức
    return t.left.push(yield) && yield(t.val) && t.right.push(yield)
}
```

## 13. Generics Best Practices

### 1. Tránh lạm dụng Generics thay thế cho Interface thông thường
Chỉ dùng generics khi kiểu dữ liệu cần giữ nguyên qua các hàm (ví dụ: Slice, Map helper, hoặc cấu trúc dữ liệu tổng quát).
```go
// ❌ Sai lầm (Lạm dụng Generics)
func Read[T io.Reader](r T)

// ✅ Đúng (Dùng Interface thông thường vì không cần bảo lưu kiểu gốc của Reader)
func Read(r io.Reader)
```

### 2. Tận dụng tối đa Type Inference
Trình biên dịch Go rất mạnh trong việc tự động suy diễn kiểu dữ liệu. Hãy để code sạch bằng cách bỏ bớt phần khai báo kiểu rõ ràng ở Client-side.
```go
// ✅ Đúng (Type inference tự hiểu)
list := NewList(1, 2, 3)

// ❌ Dư thừa
list := NewList[int](1, 2, 3)
```

## 14. Security

- **Input validation**: Validate tất cả input từ user (Sử dụng validator v10).
- **SQL injection**: Luôn dùng parameterized queries hoặc SQLC để sinh code an toàn kiểu.
- **Secrets**: Không hardcode, dùng environment variables hoặc secret manager.
- **Dependencies**: Scan vulnerabilities thường xuyên: `govulncheck ./...`.
- **TLS**: Dùng TLS cho mọi network communication.
- **Rate limiting**: Implement cho public APIs (Uber ratelimit / Leaky Bucket).

## 15. Modern Conventions & ID Generation

- **Sử dụng `any` thay cho `interface{}`**:
  - Dùng `any` cho tất cả các khai báo kiểu dữ liệu chung thay thế hoàn toàn cho `interface{}` giúp code sạch, gọn và hiện đại hơn.
- **Sinh khóa chính ID dạng UUID/ULID**:
  - Ưu tiên sử dụng **UUID v7** làm khóa chính (primary key) nhờ khả năng sắp xếp theo thời gian (time-ordered), giảm thiểu phân mảnh index trong DB.
  - Nếu không sử dụng UUID v7, sử dụng **ULID** làm giải pháp thay thế.
  - Thứ tự ưu tiên bắt buộc: **UUID v7 -> ULID**.
- **Tuyệt đối không Hardcode (No Hardcoded Literals)**:
  - Tránh hardcode trực tiếp các chuỗi định danh hệ thống (như Roles, Statuses, User Types, Permissions,...) hay các số ma thuật (magic numbers) rải rác trong mã nguồn logic xử lý.
  - Phải tập hợp và khai báo chúng một cách tường minh thành các hằng số (`const`) có ý nghĩa rõ ràng trong một package cấu trúc dùng chung (ví dụ: `internal/common`), giúp dễ dàng bảo trì và loại bỏ hoàn toàn nguy cơ gõ sai chính tả (typos).
