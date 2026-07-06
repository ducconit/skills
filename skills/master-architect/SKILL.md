---
name: master-architect
description: Master rules and guidelines for architecture, tech stack definition (Go, Vue), and Docker Compose design.
---

# Master Architect — Quy chuẩn Kiến trúc và Tech Stack Dự án

Skill này tập hợp các quy tắc thiết kế hệ thống, kiến trúc ứng dụng và định hình stack công nghệ chuẩn cho các ngôn ngữ/framework được sử dụng trong hệ thống.

Vui lòng tham khảo các tài liệu chi tiết trong thư mục `references/` bên dưới:

## 1. Quy chuẩn Docker Compose
- Hướng dẫn thiết kế cổng chuyển tiếp `FORWARD_...`, cơ chế override cục bộ qua `compose.yml` và đồng bộ file env.
- Xem chi tiết tại: [docker-compose.md](references/docker-compose.md)

## 2. Tech Stack & Quy chuẩn Golang
- Định hình bộ công cụ chuẩn (Gin, Cobra, Viper, Goose...) và các best practices khi viết code Go.
- Xem chi tiết tại: [golang.md](references/golang.md)

## 3. Tech Stack & Quy chuẩn Vue.js
- Hướng dẫn cấu trúc Vue 3 Composition API, Pinia và bắt buộc sử dụng **Bun** làm package manager.
- Xem chi tiết tại: [vue.md](references/vue.md)
