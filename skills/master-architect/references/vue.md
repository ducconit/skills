# Tech Stack & Quy chuẩn Phát triển Vue.js

Tài liệu này tập hợp các quy tắc thiết kế hệ thống, kiến trúc ứng dụng và định hình stack công nghệ chuẩn cho framework Vue.js.

## 1. Tech Stack chuẩn
- **Framework**: Vue 3 (Composition API)
- **State Management**: Pinia (Setup Store style)
- **Package Manager**: **Bun** (Bắt buộc dùng `bun` thay vì `npm/yarn/pnpm`)
- **Language**: TypeScript

## 2. Quy chuẩn Phát triển Code Vue
- **Single File Components**: Luôn sử dụng cú pháp `<script setup lang="ts">`.
- **Strict Typing**: Props/emits phải được định nghĩa kiểu dữ liệu (strongly typed) thông qua interface, tuyệt đối không dùng kiểu `any`.
- **Composables**: Trích xuất logic tái sử dụng vào các composable, luôn dọn dẹp tài nguyên (lắng nghe sự kiện, timer...) trong `onUnmounted`.
- **Hiệu năng & Tối ưu**:
  - Thực hiện lazy loading cho router và các component lớn/nặng.
  - Sử dụng `computed` cho derived state để tận dụng cache của Vue.
  - Sử dụng `shallowRef` cho các đối tượng dữ liệu lớn không cần reactive sâu (deep reactivity).
