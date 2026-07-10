---
name: golang-master
description: Comprehensive Go/Golang best practices, code review, testing, and continuous learning skill. Triggers on writing, reviewing, debugging, optimizing, or testing Go code. Specifically covers idiomatic Go patterns, recommended tech stack, performance optimization, concurrency, error handling, project structure, writing tests, unit test, integration test, mock, mockery, mock service, mock repository, assertions, table-driven tests, handler test, kiểm thử, viết test, giả lập cơ sở dữ liệu.
---

# Go Master — Kỹ năng Go Toàn diện

Skill này giúp viết Go code chất lượng cao, review code hiệu quả, và tối ưu hóa hệ thống.

## Khi nào kích hoạt skill này

- Viết hoặc chỉnh sửa code Go (API, CLI, library, microservice)
- Review hoặc refactor code Go
- Debug lỗi Go (concurrency, memory leaks, performance)
- Tối ưu hiệu năng ứng dụng Go (GC tuning, sync.Pool)
- Thiết kế package structure hoặc xử lý circular dependencies

## Hướng dẫn sử dụng

### Trước khi viết code Go
- Đọc [references/go_best_practices.md](references/go_best_practices.md) để nắm các pattern và gotchas.
- Kiểm tra [references/go_recommended_stack.md](references/go_recommended_stack.md) để dùng đúng thư viện/công cụ khuyến nghị.
- Xem [references/go_testing_best_practices.md](references/go_testing_best_practices.md) để nắm rõ cách phân tích, giả lập và viết unit test cho các layers (Handler, Service, Repository).

### Khi review hoặc tối ưu hiệu năng
- Đối chiếu với [references/go_code_review_checklist.md](references/go_code_review_checklist.md).
- Kiểm tra các kỹ thuật tối ưu hóa tại [references/go_performance_patterns.md](references/go_performance_patterns.md).

### Scripts hỗ trợ
- `scripts/go_lint_check.sh <project_dir>`: Chạy go vet, staticcheck, golangci-lint, go test -race.
- `scripts/go_complexity_check.sh <project_dir>`: Kiểm tra độ phức tạp cyclomatic (cảnh báo nếu > 15).

## Học tập liên tục

Khi phát hiện pattern Go tốt hơn hoặc bài học mới:
- Cập nhật vào [references/learned_patterns.md](references/learned_patterns.md) theo format chuẩn.
