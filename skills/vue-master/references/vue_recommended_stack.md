# Vue.js Recommended Tech Stack

> Đây là danh sách tools/library được khuyến nghị sử dụng trong dự án Vue.js.
> Agent **phải ưu tiên** dùng các tools này khi tạo project mới hoặc thêm feature.

## Package Manager

| Tool | Ghi chú |
|------|---------|
| **bun** | Package manager mặc định. Dùng `bun` thay vì npm/pnpm/yarn |

### Bun Usage

```bash
# Install dependencies
bun install

# Add dependency
bun add vue-router pinia
bun add -D vitest @vue/test-utils

# Run scripts
bun run dev
bun run build
bun run test

# Create new project
bun create vite my-project --template vue-ts
```

> **Lưu ý**: Các script trong `scripts/` vẫn dùng `npx` cho tương thích. Trong project thực tế, dùng `bunx` thay `npx`.

## Core Stack

| Mục đích | Library | Ghi chú |
|----------|---------|---------|
| **Framework** | Vue 3 | Composition API + `<script setup lang="ts">` |
| **Build Tool** | Vite | Fast HMR, native ESM |
| **Language** | TypeScript | Bắt buộc cho mọi project |

## State & Routing

| Mục đích | Library | Ghi chú |
|----------|---------|---------|
| **State Management** | Pinia | Setup stores (Composition API style) |
| **Routing** | Vue Router 4 | Lazy loading routes |

## UI & Styling

| Mục đích | Library | Ghi chú |
|----------|---------|---------|
| **CSS** | Vanilla CSS / CSS Modules | Scoped styles trong SFC |
| **Icons** | Iconify / Heroicons | Tree-shakeable |

> Component library tùy dự án (PrimeVue, Naive UI, Element Plus, etc.)

## Form & Validation

| Mục đích | Library | Ghi chú |
|----------|---------|---------|
| **Validation** | Zod hoặc Valibot | Schema-based validation |
| **Form** | VeeValidate + Zod | Nếu cần form library |

## Testing

| Mục đích | Library | Ghi chú |
|----------|---------|---------|
| **Unit/Component Test** | Vitest | Vite-native, nhanh |
| **Component Mount** | @vue/test-utils | Mount/shallow mount components |
| **E2E** | Playwright hoặc Cypress | Nếu cần E2E testing |

## Utilities

| Mục đích | Library | Ghi chú |
|----------|---------|---------|
| **Composables** | VueUse | Collection of essential composables |
| **HTTP Client** | ky hoặc ofetch | Lightweight fetch wrappers |
| **Date** | dayjs | Lightweight (2KB), NOT moment.js |

## Dev Tools

| Mục đích | Tool | Ghi chú |
|----------|------|---------|
| **Linter** | ESLint + @antfu/eslint-config | Flat config, opinionated |
| **Formatter** | Prettier | Format on save |
| **Type Check** | vue-tsc | TypeScript check cho SFC |
| **DevTools** | Vue DevTools | Browser extension |

## Summary: package.json Dependencies

Khi tạo project Vue mới, đây là baseline:

```json
{
  "dependencies": {
    "vue": "^3.5",
    "vue-router": "^4",
    "pinia": "^3"
  },
  "devDependencies": {
    "typescript": "^5",
    "vite": "^6",
    "@vitejs/plugin-vue": "^5",
    "vue-tsc": "^2",
    "vitest": "^3",
    "@vue/test-utils": "^2"
  }
}
```

> **Lưu ý**: Versions ở trên là gợi ý. Luôn kiểm tra version mới nhất khi tạo project.
> Không phải project nào cũng cần tất cả. Chỉ thêm khi thực sự cần.
