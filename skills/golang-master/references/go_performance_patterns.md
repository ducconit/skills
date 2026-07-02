# Go Performance Patterns

Các patterns tối ưu hiệu năng cho Go. Chỉ áp dụng khi đã đo và xác nhận bottleneck.

> **Nguyên tắc #1**: Đo trước, tối ưu sau. Dùng `go test -bench` và `pprof`.

## 1. String Operations

```go
// ❌ String concatenation trong loop
var s string
for _, item := range items {
    s += item.Name + ", "
}

// ✅ strings.Builder
var b strings.Builder
b.Grow(len(items) * 20) // preallocate
for i, item := range items {
    if i > 0 {
        b.WriteString(", ")
    }
    b.WriteString(item.Name)
}
result := b.String()

// ✅ strings.Join cho simple cases
result := strings.Join(names, ", ")
```

## 2. Slice Operations

```go
// ❌ Không preallocate
var result []string
for _, item := range items {
    result = append(result, item.Name)
}

// ✅ Preallocate với known capacity
result := make([]string, 0, len(items))
for _, item := range items {
    result = append(result, item.Name)
}

// ✅ Direct index khi biết chính xác size
result := make([]string, len(items))
for i, item := range items {
    result[i] = item.Name
}
```

## 3. sync.Pool & Các Cạm bẫy Trí mạng (Gotchas)

`sync.Pool` giúp giảm phân bổ bộ nhớ trên Heap cho các đối tượng tạm thời có vòng đời ngắn. Tuy nhiên, nó có các cạm bẫy sau:

1. **Ảo tưởng lưu trữ bền vững (Permanent Cache Illusion)**: `sync.Pool` KHÔNG phải là cache lưu trữ. GC có quyền xóa sạch mọi đối tượng trong `sync.Pool` bất kỳ lúc nào. Tuyệt đối không lưu các tài nguyên tĩnh như DB connection, config hay struct cần tồn tại lâu dài trong Pool.
2. **Ô nhiễm dư lượng (Stale Data)**: Đối tượng lấy từ Pool chứa dữ liệu cũ từ luồng chạy trước đó. Bạn bắt buộc phải gọi hàm dọn dẹp/reset dữ liệu đối tượng trước khi tái sử dụng hoặc trước khi đưa trở lại Pool.
3. **Rò rỉ con trỏ (Pointer Peril)**: Struct chứa con trỏ lồng ghép nếu không đặt về `nil` khi đưa lại Pool sẽ duy trì tham chiếu đến các đối tượng bên ngoài, ngăn GC thu hồi vùng nhớ đó và gây rò rỉ bộ nhớ nghiêm trọng.
4. **Sức ép phình to (Capped Pool)**: Trong môi trường tải cao, việc lưu các slice dung lượng lớn vào Pool có thể làm tăng dung lượng RAM sử dụng của app lên gấp 4 lần. Hãy giới hạn dung lượng đối tượng tối đa được phép đưa lại Pool (ví dụ: chỉ giữ lại slice < 1MB, slice lớn hơn thì để GC thu dọn tự nhiên).

```go
// Pattern chuẩn đóng gói sync.Pool an toàn
type BufferPool struct {
    pool sync.Pool
}

func NewBufferPool() *BufferPool {
    return &BufferPool{
        pool: sync.Pool{
            New: func() any {
                return new(bytes.Buffer)
            },
        },
    }
}

func (p *BufferPool) Get() *bytes.Buffer {
    return p.pool.Get().(*bytes.Buffer)
}

func (p *BufferPool) Put(buf *bytes.Buffer) {
    // Rào chắn bảo vệ sức ép phình to: Chỉ giữ lại buffer nhỏ hơn 1MB (1<<20 bytes)
    if buf.Cap() > 1024*1024 {
        return 
    }
    buf.Reset() // Đảm bảo reset stale data sạch sẽ
    p.pool.Put(buf)
}
```

## 4. Tránh interface{}/any trong Hot Paths

```go
// ❌ Generic map cho config values
config := map[string]interface{}{...}
val := config["timeout"].(int) // type assertion overhead

// ✅ Typed struct
type Config struct {
    Timeout int
    MaxRetries int
}
```

## 5. Struct Field Ordering (Memory Alignment)

```go
// ❌ Unaligned (trên 64-bit: 32 bytes do padding)
type Bad struct {
    a bool    // 1 byte + 7 padding
    b int64   // 8 bytes
    c bool    // 1 byte + 7 padding
    d int64   // 8 bytes
}

// ✅ Aligned (24 bytes, tiết kiệm 8 bytes)
type Good struct {
    b int64   // 8 bytes
    d int64   // 8 bytes
    a bool    // 1 byte
    c bool    // 1 byte + 6 padding
}

// Tool: go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
// fieldalignment -fix ./...
```

## 6. I/O với bufio

```go
// ❌ Unbuffered reads
scanner := bufio.NewScanner(file)  // ok nhưng default 64KB

// ✅ Buffered with custom size
reader := bufio.NewReaderSize(file, 256*1024) // 256KB buffer

// ✅ Buffered writes
writer := bufio.NewWriterSize(file, 256*1024)
defer writer.Flush()
```

## 7. Connection Pooling

```go
// Database connections
db, err := sql.Open("postgres", connString)
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(5 * time.Minute)
db.SetConnMaxIdleTime(1 * time.Minute)

// HTTP client
client := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}
```

## 8. Tránh Reflection trong Hot Paths

```go
// ❌ reflect trong hot path
func process(v interface{}) {
    val := reflect.ValueOf(v)
    // ... expensive reflection
}

// ✅ Type switch cho known types
func process(v interface{}) {
    switch x := v.(type) {
    case string:
        processString(x)
    case int:
        processInt(x)
    }
}

// ✅ Code generation thay cho reflection (go generate)
```

## 9. Map Operations

```go
// Preallocate maps khi biết size
m := make(map[string]int, expectedSize)

// Dùng sync.Map cho concurrent read-heavy scenarios
// Dùng sync.Mutex + map cho write-heavy scenarios
```

## 10. Điều tiết áp lực GC trong Container: GOMEMLIMIT vs GOGC

Trước Go 1.19, chỉ số `GOGC` (mặc định là 100) sẽ kích hoạt GC khi live heap tăng 100% so với chu kỳ trước. Điều này dễ dẫn đến **OOM Killed** trong môi trường container giới hạn RAM tĩnh (ví dụ: RAM 3GB, live heap 1.5GB -> GOGC cho phép phình lên 3GB -> sập trước khi kịp GC).

### Giải pháp (Go 1.19+):
Sử dụng `GOMEMLIMIT` để đặt giới hạn mềm tuyệt đối cho bộ nhớ của app.
- GC sẽ chạy dựa trên `GOGC` khi RAM còn nhiều để tiết kiệm CPU.
- Khi dung lượng vùng nhớ tiến sát ngưỡng `GOMEMLIMIT`, GC sẽ tự động gạt luật nhân đôi, tăng tần suất dọn rác tối đa để giữ RAM dưới ngưỡng giới hạn, bảo vệ app không bị OOM.
- **Best Practice trong Container**: Tắt hoàn toàn bộ hẹn giờ tương đối bằng cách cấu hình **`GOGC=off`** và thiết lập **`GOMEMLIMIT=X`** (với X bằng 90-95% giới hạn của container). Điều này tối ưu hóa CPU tối đa và chỉ chạy GC khi thực sự tiệm cận ngưỡng sập RAM.

## 11. Khử chi phí cấp phát bằng Zero-Copy Conversion

Chuyển đổi chuỗi ký tự (`string` - bất biến) sang mảng bytes (`[]byte` - có thể thay đổi) thông thường buộc Go phải cấp phát bộ nhớ mới và sao chép dữ liệu.

### Tùy chọn 1 (Go 1.20+ unsafe):
Sử dụng gói `unsafe` để chỉ trỏ Slice Header vào cùng vùng bộ nhớ gốc mà không cần sao chép (tối ưu hóa zero-allocation, giảm từ ~60ns xuống 0.5ns).
```go
import "unsafe"

// Chuyển string sang []byte không allocation
func StringToBytes(s string) []byte {
    return unsafe.Slice(unsafe.StringData(s), len(s))
}

// Chuyển []byte sang string không allocation
func BytesToString(b []byte) string {
    return unsafe.String(unsafe.SliceData(b), len(b))
}
```

### Tùy chọn 2 (Tối ưu hóa tự động của trình biên dịch Go 1.22+):
Nếu biến mảng bytes sinh ra từ chuỗi chỉ dùng để **đọc (read-only usage)** và không bị sửa đổi trong suốt vòng đời, trình biên dịch Go 1.22+ sẽ tự động biên dịch bằng cơ chế tối ưu PCDATA (Zero-copy) ngầm. Không cần thiết phải viết code `unsafe` thủ công.

## 12. In-Memory Cache hiệu năng cao loại bỏ áp lực GC (Zero GC Cache)

Cơ chế thu gom rác của Go quét đồ thị đối tượng (Mark-and-Sweep). Nếu bạn dùng cache chứa hàng triệu con trỏ, GC pause sẽ tăng vọt (lên hàng trăm ms) do phải quét qua đồ thị con trỏ khổng lồ.

### Nguyên lý tối ưu (Go 1.5+):
*"Bất kỳ map nào không chứa con trỏ trong cả key lẫn value (ví dụ `map[uint64]uint32`), GC sẽ bỏ qua hoàn toàn quá trình quét cấu trúc map đó (blind array)."*

### Thiết kế theo BigCache / FreeCache:
1. Băm key dạng chuỗi thành số nguyên `uint64` (dùng làm key của map).
2. Lưu dữ liệu thô (JSON bytes) nối tiếp nhau vào một mảng bytes lớn được cấp phát trước (`ring-buffer []byte`).
3. Lưu vị trí offset (`uint32`) của dữ liệu trong ring-buffer vào map làm value.
4. Đồ thị map có dạng `map[uint64]uint32` (hoàn toàn không chứa con trỏ), giúp GC chạy với tốc độ cực nhanh ngay cả khi cache lưu hàng GB dữ liệu.

## 13. Benchmarking

```go
func BenchmarkProcess(b *testing.B) {
    data := setupTestData()
    b.ResetTimer() // reset sau setup

    for i := 0; i < b.N; i++ {
        process(data)
    }
}

// Run: go test -bench=BenchmarkProcess -benchmem -count=5
// Compare: benchstat old.txt new.txt
```

## 14. Profiling

```bash
# CPU profile
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profile
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof

# Trace
go test -trace=trace.out -bench=.
go tool trace trace.out

# In production (net/http/pprof)
import _ "net/http/pprof"
go func() { http.ListenAndServe(":6060", nil) }()
# curl http://localhost:6060/debug/pprof/heap > heap.prof
```
