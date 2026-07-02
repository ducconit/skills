---
name: vue-master
description: Comprehensive Vue.js 3 best practices, code review, and continuous learning skill. Use when writing, reviewing, debugging, or optimizing Vue.js applications. Covers Composition API, state management (Pinia), routing, TypeScript integration, component design, performance optimization, and testing.
---

# Vue Master — Kỹ năng Vue.js 3 Toàn diện

Skill này giúp agent viết Vue.js code chất lượng cao, review code hiệu quả, và liên tục cải thiện kiến thức Vue/frontend.

## Khi nào kích hoạt skill này

- Viết Vue.js components, composables, stores mới
- Review hoặc refactor Vue.js code
- Debug issues liên quan đến reactivity, performance, rendering
- Tối ưu bundle size và performance
- Thiết kế component architecture
- Viết tests cho Vue.js applications
- Tích hợp TypeScript vào Vue project

## Hướng dẫn sử dụng

### Trước khi viết code Vue

1. Đọc `references/vue_best_practices.md` để nắm các pattern chuẩn
2. Kiểm tra `examples/` để xem code mẫu idiomatic
3. Kiểm tra `references/vue_recommended_stack.md` để dùng đúng tools/library
4. Áp dụng các nguyên tắc cốt lõi:
   - **Composition API over Options API**: Luôn dùng `<script setup lang="ts">`
   - **TypeScript first**: Type everything, avoid `any`
   - **Reactive by design**: Hiểu rõ ref vs reactive, computed vs watch
   - **Small components**: Mỗi component < 200 lines, một trách nhiệm
   - **Composables for reuse**: Extract shared logic vào composables
   - **Bun as package manager**: Dùng `bun` thay vì npm/pnpm/yarn

### Khi review code Vue

1. Chạy script `scripts/vue_lint_check.sh` để kiểm tra tự động
2. Chạy `scripts/vue_bundle_check.sh` để kiểm tra bundle size
3. Đối chiếu với `references/vue_code_review_checklist.md`
4. Kiểm tra performance patterns tại `references/vue_performance_patterns.md`

### Scripts có sẵn

| Script | Mục đích |
|--------|----------|
| `scripts/vue_lint_check.sh <project_dir>` | Chạy ESLint, vue-tsc, Prettier checks |
| `scripts/vue_bundle_check.sh <project_dir>` | Phân tích bundle size, phát hiện dependencies lớn |

### References có sẵn

| File | Nội dung |
|------|----------|
| `references/vue_best_practices.md` | Best practices toàn diện cho Vue.js 3 |
| `references/vue_recommended_stack.md` | Tech stack và tools được khuyến nghị |
| `references/vue_code_review_checklist.md` | Checklist review code Vue |
| `references/vue_performance_patterns.md` | Patterns tối ưu hiệu năng Vue |
| `references/learned_patterns.md` | Patterns đã học được (liên tục cập nhật) |

### Code mẫu

| File | Pattern |
|------|---------|
| `examples/composable_pattern.ts` | Composable pattern (useFetch) |
| `examples/typed_component.vue` | Full TypeScript SFC component |
| `examples/pinia_store.ts` | Pinia store với Composition API |
| `examples/route_setup.ts` | Vue Router lazy loading & guards |

## Học tập liên tục

Khi phát hiện pattern Vue mới tốt hơn, best practice mới, hoặc anti-pattern cần tránh:

1. **Cập nhật** `references/learned_patterns.md` với entry mới theo format:
   ```
   ### [YYYY-MM-DD] Tiêu đề Pattern
   **Context**: Tình huống phát hiện
   **Pattern**: Mô tả pattern
   **Lý do tốt hơn**: Giải thích
   **Ví dụ**: Code minh họa
   ```

2. Nếu pattern đủ quan trọng, **thêm vào** `references/vue_best_practices.md` ở section phù hợp

3. Nếu có code mẫu hay, **thêm file mới** vào `examples/`

4. Nếu cần tool kiểm tra mới, **thêm script** vào `scripts/`
