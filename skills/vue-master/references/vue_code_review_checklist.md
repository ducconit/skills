# Vue.js Code Review Checklist

Checklist dùng khi review code Vue.js. Đảm bảo tất cả các mục đều pass trước khi approve.

## Script & TypeScript

- [ ] Dùng `<script setup lang="ts">` (Composition API + TypeScript)
- [ ] Props typed với interface/type, có defaults nếu optional
- [ ] Emits typed đầy đủ với payload types
- [ ] Không dùng `any` type (dùng `unknown` nếu thực sự cần)
- [ ] Template refs typed chính xác (`ref<HTMLInputElement | null>(null)`)

## Reactivity

- [ ] `ref` cho primitives, `reactive` cho objects (hoặc `ref` cho cả hai)
- [ ] `computed` cho derived state (không dùng methods cho cached values)
- [ ] Watchers có cleanup logic (timers, event listeners, subscriptions)
- [ ] `shallowRef` cho large objects không cần deep reactivity
- [ ] Không mutate props trực tiếp

## Component Design

- [ ] Component < 200 lines (tách nếu quá dài)
- [ ] Single responsibility — mỗi component làm 1 việc
- [ ] Reusable logic extracted vào composables
- [ ] Không prop drilling (dùng provide/inject hoặc Pinia cho 3+ levels)
- [ ] Slots cho flexible content

## Template

- [ ] `key` attribute cho tất cả `v-for` loops (dùng unique ID, không dùng index)
- [ ] `v-if` và `v-for` không cùng element (dùng `<template>` wrapper)
- [ ] Event handlers named convention: `@click="onClickXxx"`
- [ ] Không logic phức tạp trong template (move vào computed/methods)
- [ ] `data-testid` cho elements cần test

## Styling

- [ ] CSS `scoped` hoặc CSS Modules (tránh global styles trong components)
- [ ] Không inline styles (trừ dynamic styles từ props)
- [ ] BEM hoặc naming convention nhất quán
- [ ] Responsive design considerations

## State Management

- [ ] Pinia stores dùng Setup Store syntax (Composition API)
- [ ] Store actions cho async operations
- [ ] Không access store directly từ templates (dùng composables nếu cần logic)
- [ ] Store state properly typed

## Performance

- [ ] Routes lazy loaded (`() => import(...)`)
- [ ] Heavy components dùng `defineAsyncComponent`
- [ ] `v-once` cho static content
- [ ] `KeepAlive` cho components hay switch qua lại
- [ ] Không re-render không cần thiết (check computed vs watch)

## Memory & Cleanup

- [ ] Event listeners removed trong `onUnmounted`
- [ ] Timers (setTimeout, setInterval) cleared
- [ ] Subscriptions unsubscribed
- [ ] Abort controllers cho fetch requests

## Security

- [ ] Không `v-html` với user input
- [ ] Input sanitized/validated
- [ ] Không expose secrets trong client code
- [ ] API calls qua HTTPS

## Accessibility

- [ ] Semantic HTML elements (button, nav, main, article...)
- [ ] ARIA labels cho interactive elements
- [ ] Keyboard navigation support
- [ ] Focus management cho modals/dialogs
