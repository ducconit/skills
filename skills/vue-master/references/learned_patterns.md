# Learned Vue.js Patterns & Insights

> File này được liên tục cập nhật khi phát hiện patterns Vue.js mới, best practices, hoặc anti-patterns cần tránh.
> Mỗi entry ghi lại context, pattern, lý do, và ví dụ cụ thể.

---

### [2026-06-29] Format Entry Template

**Context**: Đây là template cho các entries mới.

**Pattern**: Mô tả pattern/practice đã học.

**Lý do tốt hơn**: Giải thích tại sao pattern này tốt hơn cách cũ.

**Ví dụ**:
```vue
<!-- Code minh họa pattern -->
```

---

### [2026-06-29] Dùng `defineModel` thay v-model manual (Vue 3.4+)

**Context**: Trước Vue 3.4, custom v-model cần khai báo props + emit thủ công.

**Pattern**: `defineModel()` tạo two-way binding trực tiếp, giảm boilerplate.

**Lý do tốt hơn**: Ít code hơn, type-safe, dễ đọc hơn.

**Ví dụ**:
```vue
<!-- Trước Vue 3.4 -->
<script setup lang="ts">
const props = defineProps<{ modelValue: string }>()
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()
</script>
<template>
  <input :value="props.modelValue" @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)" />
</template>

<!-- Vue 3.4+ với defineModel -->
<script setup lang="ts">
const model = defineModel<string>()
</script>
<template>
  <input v-model="model" />
</template>
```

---

### [2026-06-29] `toValue()` thay `unref()` trong composables

**Context**: Composables nên accept cả ref, getter, và plain values.

**Pattern**: Dùng `toValue()` (Vue 3.3+) thay `unref()` vì nó cũng handle getter functions.

**Lý do tốt hơn**: Linh hoạt hơn, cho phép truyền `() => someValue` ngoài ref và plain value.

**Ví dụ**:
```typescript
import { toValue, type MaybeRefOrGetter } from 'vue'

export function useFetch(url: MaybeRefOrGetter<string>) {
  watchEffect(() => {
    fetch(toValue(url)) // works with ref, getter, or plain string
  })
}

// Tất cả đều hợp lệ:
useFetch('/api/users')
useFetch(urlRef)
useFetch(() => `/api/users/${userId.value}`)
```
