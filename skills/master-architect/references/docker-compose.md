# Quy chuẩn Thiết kế và Kiến trúc Docker Compose

Tài liệu này định hình cách thiết kế, viết và quản lý Docker Compose của dự án một cách chuyên nghiệp, sạch sẽ và dễ dàng mở rộng.

## Nguyên tắc Thiết kế Cốt lõi

1. **Độc lập và Mặc định Hợp lý (Standalone & Fallbacks)**:
   - Không được hardcode port máy chủ (host port) hoặc mật khẩu trực tiếp trong file cấu hình gốc.
   - Luôn sử dụng cú pháp `${FORWARD_...:-default_value}` để lấy cấu hình từ biến môi trường.
   - Cung cấp giá trị fallback mặc định (`default_value`) là cổng gốc/chuẩn của dịch vụ đó, giúp Docker Compose vẫn khởi chạy thành công ngay cả khi không có file `.env`.

2. **Quy chuẩn Đặt tên Cổng chuyển tiếp (FORWARD_ prefix)**:
   - Toàn bộ cổng host ánh xạ ra ngoài (host ports - phần bên trái dấu `:`) phải sử dụng tiền tố `FORWARD_` (ví dụ: `FORWARD_DB_PRIMARY_PORT`, `FORWARD_REDIS_PORT`, `FORWARD_MINIO_PORT`).
   - Việc này giúp phân biệt rõ ràng giữa biến cấu hình cổng của Docker Compose (`FORWARD_...`) và biến cấu hình kết nối của chính ứng dụng Backend (như `REDIS_PORT`, `MAIL_PORT`...).

3. **Cơ chế Ghi đè cục bộ (Override Mechanism)**:
   - File Docker Compose chính thống của dự án được đặt tên là `docker-compose.yml`.
   - Các tùy chỉnh mang tính cá nhân, cục bộ của từng lập trình viên (developer) sẽ được viết trong file `compose.yml`.
   - File `compose.yml` **phải** được đưa vào `.gitignore` để tránh đẩy lên repository chung.

4. **Đồng bộ hóa Tài liệu Cấu hình (.env.example)**:
   - Mọi biến cấu hình động mới được thêm vào Docker Compose đều phải được khai báo làm mẫu trong file `.env.example` kèm theo các ghi chú cụ thể để lập trình viên khác biết và sử dụng.

---

## Ví dụ Thực tế

### 1. Cấu hình Dịch vụ Cơ sở dữ liệu (`docker-compose.yml`)
```yaml
services:
  postgres:
    image: postgres:18-alpine
    container_name: myapp-postgres
    environment:
      POSTGRES_USER: ${DB_PRIMARY_USER:-postgres}
      POSTGRES_PASSWORD: ${DB_PRIMARY_PASSWORD:-postgres}
      POSTGRES_DB: ${DB_PRIMARY_NAME:-myapp}
    ports:
      - "${FORWARD_DB_PRIMARY_PORT:-5432}:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_PRIMARY_USER:-postgres}"]
      interval: 5s
      timeout: 5s
      retries: 5
```

### 2. Thiết lập Biến Môi trường tương ứng (`.env.example`)
```bash
# App Database Connection
DB_PRIMARY_HOST=127.0.0.1
DB_PRIMARY_PORT=5432
DB_PRIMARY_USER=postgres
DB_PRIMARY_PASSWORD=postgres
DB_PRIMARY_NAME=myapp

# Docker Compose Services Port Forwarding (Host Ports)
FORWARD_DB_PRIMARY_PORT=5432
```
