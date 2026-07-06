# Tech Stack & Quy chuẩn Phát triển Golang

Tài liệu này tập hợp các quy tắc thiết kế hệ thống, kiến trúc ứng dụng và định hình stack công nghệ chuẩn cho ngôn ngữ Golang.

## 1. Tech Stack chuẩn
- **API Framework**: Gin
- **CLI & Configurations**: Cobra, Viper
- **Database Migration**: Goose
- **Logger**: zerolog
- **Validation**: validator v10
- **Toolkit**: gotoolkit

## 2. Quy chuẩn Phát triển Code Go
- **Explicit Error Handling**: Luôn kiểm tra lỗi trực tiếp, wrap error kèm theo context rõ ràng.
- **Accept Interfaces, Return Structs**: Thiết kế hàm nhận interface và trả về struct cụ thể để tối đa hóa tính linh hoạt.
- **Table-Driven Tests**: Viết kiểm thử dạng bảng (table-driven), tận dụng `t.Helper()` và `t.Cleanup()` để dọn dẹp tài nguyên.
- **Tối ưu hóa**: Tránh tối ưu hóa sớm (premature optimization). Sử dụng benchmark và pprof để đo lường trước khi quyết định tối ưu.
- **Concurrency**: Luôn có cơ chế dừng goroutine, luôn truyền và sử dụng `context.Context` để quản lý vòng đời và timeout.
