# Vue.js Performance Patterns

Các patterns tối ưu hiệu năng cho Vue.js 3. Chỉ áp dụng khi đã đo và xác nhận bottleneck.

> **Nguyên tắc #1**: Đo trước, tối ưu sau. Dùng Vue DevTools Performance tab và Lighthouse.

## 1. Component Lazy Loading

```typescript
// Route-level code splitting
const routes = [
  {
    path: '/dashboard',
    component: () => import('./views/Dashboard.vue'),
  },
]

// Component-level lazy loading
import { defineAsyncComponent } from 'vue'

const HeavyChart = defineAsyncComponent(() =>
  import('./components/HeavyChart.vue')
)

// With loading/error states
const HeavyChart = defineAsyncComponent({
  loader: () => import('./components/HeavyChart.vue'),
  loadingComponent: LoadingSpinner,
  errorComponent: ErrorDisplay,
  delay: 200,        // delay before showing loading
  timeout: 10000,    // timeout
})
```

## 2. Virtual Scrolling cho Large Lists

```vue
<!-- Dùng @tanstack/vue-virtual hoặc vue-virtual-scroller -->
<script setup lang="ts">
import { useVirtualizer } from '@tanstack/vue-virtual'

const parentRef = ref<HTMLElement | null>(null)
const virtualizer = useVirtualizer({
  count: 10000,
  getScrollElement: () => parentRef.value,
  estimateSize: () => 50,
})
</script>

<template>
  <div ref="parentRef" style="height: 400px; overflow: auto">
    <div :style="{ height: `${virtualizer.getTotalSize()}px`, position: 'relative' }">
      <div
        v-for="item in virtualizer.getVirtualItems()"
        :key="item.key"
        :style="{ position: 'absolute', top: `${item.start}px`, height: `${item.size}px` }"
      >
        {{ items[item.index].name }}
      </div>
    </div>
  </div>
</template>
```

## 3. v-once cho Static Content

```vue
<!-- ✅ Content không bao giờ thay đổi -->
<template>
  <header v-once>
    <h1>Welcome to My App</h1>
    <p>This is a static description that never changes.</p>
  </header>

  <!-- Dynamic content bên dưới -->
  <main>
    <UserList :users="users" />
  </main>
</template>
```

## 4. v-memo cho Expensive Re-renders

```vue
<!-- ✅ Chỉ re-render khi dependencies thay đổi -->
<template>
  <div v-for="item in list" :key="item.id" v-memo="[item.id, item.selected]">
    <ExpensiveComponent :item="item" />
  </div>
</template>
```

## 5. KeepAlive cho Cached Components

```vue
<!-- ✅ Cache component state khi switch tabs -->
<template>
  <KeepAlive :max="5" :include="['UserProfile', 'Dashboard']">
    <component :is="currentTab" />
  </KeepAlive>
</template>

<script setup lang="ts">
// Lifecycle hooks riêng cho KeepAlive
import { onActivated, onDeactivated } from 'vue'

onActivated(() => {
  // Component được reactivate từ cache
  refreshData()
})

onDeactivated(() => {
  // Component bị deactivate (nhưng vẫn cached)
  clearTimers()
})
</script>
```

## 6. shallowRef/shallowReactive cho Large Objects

```typescript
// ❌ Deep reactivity cho large array (expensive)
const items = ref<Item[]>(largeArray)

// ✅ Shallow reactivity — chỉ track reference change
const items = shallowRef<Item[]>(largeArray)

// Phải replace toàn bộ ref để trigger update
function addItem(item: Item) {
  items.value = [...items.value, item] // ✅ triggers update
  // items.value.push(item) // ❌ KHÔNG trigger update
}

// triggerRef cho manual trigger
import { triggerRef } from 'vue'
items.value.push(newItem)
triggerRef(items) // ✅ Force trigger
```

## 7. Debouncing Watchers & Event Handlers

```typescript
import { watchDebounced } from '@vueuse/core'

// ✅ Debounced watcher
watchDebounced(
  searchQuery,
  (query) => fetchResults(query),
  { debounce: 300, maxWait: 1000 }
)

// ✅ Debounced event handler
import { useDebounceFn } from '@vueuse/core'

const debouncedSearch = useDebounceFn((query: string) => {
  fetchResults(query)
}, 300)
```

## 8. Image Optimization

```vue
<template>
  <!-- ✅ Lazy loading images -->
  <img
    :src="imageUrl"
    loading="lazy"
    :alt="altText"
    :width="300"
    :height="200"
  />

  <!-- ✅ Responsive images -->
  <picture>
    <source :srcset="imageWebP" type="image/webp" />
    <source :srcset="imageJpg" type="image/jpeg" />
    <img :src="imageFallback" :alt="altText" loading="lazy" />
  </picture>
</template>
```

## 9. Tree-Shaking Friendly Imports

```typescript
// ❌ Import toàn bộ library
import * as lodash from 'lodash'
lodash.debounce(...)

// ✅ Named imports (tree-shakeable)
import { debounce } from 'lodash-es'

// ❌ Import toàn bộ icon library
import * as Icons from '@heroicons/vue/24/solid'

// ✅ Import specific icons
import { ChevronDownIcon } from '@heroicons/vue/24/solid'
```

## 10. Web Workers cho Heavy Computation

```typescript
// worker.ts
self.addEventListener('message', (e) => {
  const result = heavyComputation(e.data)
  self.postMessage(result)
})

// composable
export function useWorker() {
  const worker = new Worker(new URL('./worker.ts', import.meta.url), {
    type: 'module',
  })

  const result = ref(null)

  worker.addEventListener('message', (e) => {
    result.value = e.data
  })

  function compute(data: unknown) {
    worker.postMessage(data)
  }

  onUnmounted(() => worker.terminate())

  return { result, compute }
}
```

## 11. Computed vs Methods

```vue
<script setup lang="ts">
// ✅ computed: cached, chỉ re-evaluate khi dependencies thay đổi
const sortedItems = computed(() =>
  [...items.value].sort((a, b) => a.name.localeCompare(b.name))
)

// ❌ method: re-evaluate mỗi lần render
function getSortedItems() {
  return [...items.value].sort((a, b) => a.name.localeCompare(b.name))
}
</script>

<template>
  <!-- ✅ Dùng computed -->
  <div v-for="item in sortedItems" :key="item.id">{{ item.name }}</div>

  <!-- ❌ Dùng method (re-sort mỗi render) -->
  <div v-for="item in getSortedItems()" :key="item.id">{{ item.name }}</div>
</template>
```

## 12. Conditional Rendering Performance

```vue
<template>
  <!-- v-show: toggle CSS display (giữ DOM), tốt cho toggle thường xuyên -->
  <HeavyComponent v-show="isVisible" />

  <!-- v-if: mount/unmount DOM, tốt cho ít toggle -->
  <HeavyComponent v-if="isVisible" />

  <!-- Rule of thumb:
       - Toggle thường xuyên → v-show
       - Toggle ít hoặc initial false → v-if
  -->
</template>
```
