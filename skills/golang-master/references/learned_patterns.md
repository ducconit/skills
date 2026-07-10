# Learned Go Patterns & Insights

> File này được liên tục cập nhật khi phát hiện patterns Go mới, best practices, hoặc anti-patterns cần tránh.
> Mỗi entry ghi lại context, pattern, lý do, và ví dụ cụ thể.
> Khi entry đủ quan trọng, nó sẽ được tích hợp vào `go_best_practices.md` hoặc `go_performance_patterns.md` rồi xóa khỏi đây.

---

### [2026-07-06] Idiomatic Go Factory: Return Structs, Accept Interfaces (Ví dụ qua Storage Manager)

**Context**: Khi viết các factory function / constructor (như `NewManager`) khởi tạo tài nguyên cho package (như `storage` hay `cache`), việc trả về interface trực tiếp từ constructor là một Go anti-pattern. Thay vào đó, cần trả về một struct cụ thể (hoặc con trỏ struct) và để phía client tự quyết định dùng interface hoặc mock.

**Pattern**:
- **Constructor trả về Struct**: Thiết kế `func NewManager(cfg config.StorageConfig) *Manager` trả về `*Manager` (struct cụ thể), thay vì trả về `Storage` (interface).
- **Client nhận Interface**: Các lớp nghiệp vụ tiêu thụ (như `fileService`) nhận tham số dạng interface `storage.Storage` để dễ dàng viết Mock test cho chính lớp đó.

**Lý do tốt hơn**:
- Tuân thủ triệt để tôn chỉ Go: *"Accept interfaces, return structs"*. Trả về struct cụ thể giúp client dễ dàng mở rộng và sử dụng đầy đủ các API không có trong interface nguyên bản mà không bị ép buộc (locked-in).
- Giúp các công cụ tự sinh mã (như Google Wire) và IDE biên dịch hoạt động hiệu quả nhất, tránh rườm rà khi Binding kiểu dữ liệu interface lúc khởi tạo DI.

**Ví dụ**:
```go
// storage/storage.go
type Storage interface {
	Disk(name string) (*blob.Bucket, error)
}

type Manager struct {
	// fields
}

// Trả về struct cụ thể *Manager
func NewManager(cfg config.StorageConfig) *Manager {
	return &Manager{...}
}

// internal/service/file.go
type fileService struct {
	storageMgr storage.Storage // Nhận interface để dễ Mock
}
```

---

### [2026-07-06] Laravel-style Multi-store Cache Abstraction với gocache v4 và Ristretto v2

**Context**: Refactor package `cache` để hỗ trợ đổi driver linh hoạt tùy theo môi trường (`redis`, `memory`, `memcached`), hỗ trợ cấu hình riêng biệt cho từng store kèm cơ chế fallback về cấu hình global và thiết lập prefix an toàn tránh trùng lặp key giữa các dự án.

**Pattern**:
1. **Graceful Cleanup Callback**: Bọc `gocache.Cache` cùng một callback dọn dẹp `onClose func() error` để đóng kết nối của các driver cụ thể (như Redis/Ristretto) khi ứng dụng tắt. Do gocache không quản lý vòng đời connection, pattern này giúp cô lập logic thao tác cache mà vẫn dọn dẹp được tài nguyên.
2. **Generics Ristretto v2 Integration**: Khi tích hợp Ristretto v2 với `gocache/store/ristretto/v4`, cần chỉ định rõ generic type parameters cho cả `Config` và `NewCache` để đồng bộ kiểu Key-Value: `ristretto.NewCache[string, any](&ristretto.Config[string, any]{...})`.
3. **Laravel-style Stores Config with Fallback**: Thiết kế cấu hình cache có mục `stores` để ghi đè các cấu hình kết nối cụ thể cho từng store, nếu thiếu sẽ tự động fallback về cấu hình global.

**Lý do tốt hơn**:
- Cho phép đổi driver cache chỉ bằng cấu hình `.env`/`config.yml` mà không phải thay đổi bất kỳ dòng code nghiệp vụ nào.
- Tránh rò rỉ socket/connection pool nhờ cơ chế dọn dẹp callback kín kẽ, đồng thời xoá bỏ hoàn toàn dependencies không dùng tới (`go-cache` cũ).

**Ví dụ**:
```go
type Cache struct {
	gocache *cache.Cache[any]
	prefix  string
	onClose func() error
}

func New(cfg config.RedisConfig, cacheCfg config.CacheConfig) (*Cache, error) {
	var gocache *cache.Cache[any]
	var onClose func() error

	driver := cacheCfg.Driver
	if driver == "" {
		driver = "memory"
	}

	switch driver {
	case "memory", "ristretto":
		client, _ := ristretto.NewCache[string, any](&ristretto.Config[string, any]{
			NumCounters: 1e7,
			MaxCost:     1e6,
			BufferItems: 64,
		})
		ristStore := ristretto_store.NewRistretto[string, any](client)
		gocache = cache.New[any](ristStore)
		onClose = func() error {
			client.Close()
			return nil
		}
	// các driver khác...
	}

	return &Cache{
		gocache: gocache,
		prefix:  cacheCfg.Prefix,
		onClose: onClose,
	}, nil
}

func (c *Cache) Close() error {
	if c.onClose != nil {
		return c.onClose()
	}
	return nil
}
```

### [2026-07-07] Layered Mocking, Gin Context Status Assertions & JSON Field Privacy Filtering

**Context**: Khi viết unit test cho HTTP Handler layer sử dụng Gin và Mockery, việc assert status code trực tiếp trên `httptest.ResponseRecorder` của Gin có thể trả về `200` thay vì `204` do dữ liệu chưa được flush. Đồng thời, cần tránh lỗi Import Cycle khi mock nhiều layer và đảm bảo bảo mật dữ liệu JSON đầu ra.

**Pattern**:
1. **Gin Context Status Assertion**: Sử dụng `c.Writer.Status()` thay vì `w.Code` để assert HTTP status code khi gọi trực tiếp handler (không đi qua router engine).
2. **Subpackage Mocks scoped layout**: Sinh mock của package nào thì nằm ở subpackage `mocks` của package đó (ví dụ: `internal/repository/mocks`, `internal/service/mocks`), tránh gom chung vào một package `mocks` duy nhất ở root gây lỗi vòng lặp import.
3. **Response Privacy Filtering**: Unmarshal response data sang `map[string]any` thô, assert sự tồn tại của các trường an toàn và sự **vắng mặt** (`nil`) của các trường nhạy cảm như `password`, `password_hash`, `secret`.

**Lý do tốt hơn**:
- Ngăn chặn lỗi kiểm thử sai lệch khi test các response không chứa body (như 204 NoContent).
- Bẻ gãy hoàn toàn lỗi Import Cycle trong Go test compilation.
- Kiểm thử bảo mật trực tiếp ở mức độ JSON response thực tế trả về cho client.

**Ví dụ**:
```go
// Setup Gin test context
w := httptest.NewRecorder()
c, _ := gin.CreateTestContext(w)

// Gọi handler trực tiếp
handler.Delete(c)

// Assert chính xác bằng Writer.Status()
assert.Equal(t, http.StatusNoContent, c.Writer.Status())

// Assert bảo mật
var resp res.Response
json.Unmarshal(w.Body.Bytes(), &resp)
dataMap := resp.Data.(map[string]any)
assert.Nil(t, dataMap["password"])
```

---

**Context**: Tình huống phát hiện.

**Pattern**: Mô tả pattern/practice đã học.

**Lý do tốt hơn**: Giải thích tại sao pattern này tốt hơn cách cũ.

**Ví dụ**:
```go
// Code minh họa pattern
```
