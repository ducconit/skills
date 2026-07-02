# Vue.js 3 Best Practices — Tham chiếu Toàn diện

## 1. Project Structure

### Feature-based Layout (Recommended)

```
src/
├── assets/                 # Static assets
├── components/             # Shared/global components
│   ├── ui/                 # Base UI components (Button, Input, Modal)
│   └── layout/             # Layout components (Header, Sidebar)
├── composables/            # Shared composables (useAuth, useFetch)
├── features/               # Feature modules
│   ├── auth/
│   │   ├── components/     # Feature-specific components
│   │   ├── composables/    # Feature-specific composables
│   │   ├── stores/         # Feature-specific stores
│   │   ├── types/          # Feature-specific types
│   │   └── views/          # Feature pages
│   └── dashboard/
├── router/                 # Route definitions
├── stores/                 # Global Pinia stores
├── types/                  # Shared TypeScript types
├── utils/                  # Utility functions
├── App.vue
└── main.ts
```

### Nguyên tắc

- Group by feature, không group by type
- Components con nằm trong feature folder
- Shared components ở `components/`
- Mỗi feature tự chứa composables, stores, types riêng

## 2. Composition API Patterns

### Script Setup (Bắt buộc)

```vue
<!-- ✅ Luôn dùng script setup + TypeScript -->
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { User } from '@/types'

// Props
const props = defineProps<{
  userId: string
  showAvatar?: boolean
}>()

// Emits
const emit = defineEmits<{
  update: [user: User]
  delete: [id: string]
}>()

// Reactive state
const user = ref<User | null>(null)
const loading = ref(false)

// Computed
const displayName = computed(() => {
  return user.value?.name ?? 'Anonymous'
})

// Lifecycle
onMounted(async () => {
  loading.value = true
  user.value = await fetchUser(props.userId)
  loading.value = false
})
</script>
```

### Composables

```typescript
// ✅ Composable conventions
// - Tên bắt đầu bằng "use"
// - Return reactive values
// - Cleanup trong onUnmounted

export function useFetch<T>(url: MaybeRef<string>) {
  const data = ref<T | null>(null)
  const error = ref<Error | null>(null)
  const loading = ref(false)

  async function execute() {
    loading.value = true
    error.value = null
    try {
      const response = await fetch(unref(url))
      data.value = await response.json()
    } catch (e) {
      error.value = e as Error
    } finally {
      loading.value = false
    }
  }

  onMounted(execute)

  return { data, error, loading, refresh: execute }
}
```

## 3. Component Design

### Nguyên tắc

- **Single Responsibility**: Mỗi component làm 1 việc
- **< 200 lines**: Nếu dài hơn, tách component
- **Props down, events up**: Data flow rõ ràng
- **Không prop drilling**: Dùng provide/inject hoặc Pinia

### Props & Emits Typing

```vue
<script setup lang="ts">
// ✅ Interface-based props
interface Props {
  title: string
  count?: number
  items: string[]
  status: 'active' | 'inactive'
}

const props = withDefaults(defineProps<Props>(), {
  count: 0,
  status: 'active',
})

// ✅ Typed emits
const emit = defineEmits<{
  'update:count': [value: number]
  submit: [data: FormData]
  cancel: []
}>()
</script>
```

### Slots

```vue
<!-- Parent -->
<DataTable :items="users">
  <template #header>
    <h2>Users</h2>
  </template>
  <template #row="{ item }">
    <td>{{ item.name }}</td>
  </template>
  <template #empty>
    <p>No users found</p>
  </template>
</DataTable>
```

## 4. State Management with Pinia

### Setup Store (Recommended)

```typescript
// ✅ Composition API style store
export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)

  // Getters (computed)
  const isAuthenticated = computed(() => !!token.value)
  const displayName = computed(() => user.value?.name ?? 'Guest')

  // Actions
  async function login(credentials: LoginCredentials) {
    const response = await authApi.login(credentials)
    user.value = response.user
    token.value = response.token
  }

  function logout() {
    user.value = null
    token.value = null
  }

  return { user, token, isAuthenticated, displayName, login, logout }
})
```

### Store Composition

```typescript
export const useCartStore = defineStore('cart', () => {
  const authStore = useAuthStore() // Dùng store khác
  const items = ref<CartItem[]>([])

  const total = computed(() =>
    items.value.reduce((sum, item) => sum + item.price * item.quantity, 0)
  )

  async function checkout() {
    if (!authStore.isAuthenticated) {
      throw new Error('Must be authenticated')
    }
    // ...
  }

  return { items, total, checkout }
})
```

## 5. Vue Router

### Lazy Loading Routes

```typescript
const routes = [
  {
    path: '/',
    component: () => import('@/features/home/views/HomePage.vue'),
  },
  {
    path: '/dashboard',
    component: () => import('@/features/dashboard/views/DashboardLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        component: () => import('@/features/dashboard/views/DashboardHome.vue'),
      },
    ],
  },
]
```

### Navigation Guards

```typescript
router.beforeEach(async (to) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
})
```

## 6. TypeScript Integration

- Luôn dùng `<script setup lang="ts">`
- Type props với interface, không dùng runtime validation trừ khi cần
- Type emits đầy đủ
- Type template refs: `const inputRef = ref<HTMLInputElement | null>(null)`
- Type composable return values
- Dùng `satisfies` operator cho type checking

## 7. Reactivity System

```typescript
// ref cho primitives
const count = ref(0)
const name = ref('')

// reactive cho objects (nhưng ref cũng ok)
const form = reactive({
  name: '',
  email: '',
})

// computed cho derived state (cached, lazy)
const fullName = computed(() => `${first.value} ${last.value}`)

// watch cho side effects
watch(searchQuery, (newQuery) => {
  fetchResults(newQuery)
}, { debounce: 300 })

// watchEffect cho auto-tracking
watchEffect(() => {
  console.log(`Count is: ${count.value}`)
})

// shallowRef cho large objects (performance)
const largeList = shallowRef<Item[]>([])
// Phải replace toàn bộ ref để trigger reactivity
largeList.value = [...largeList.value, newItem]
```

## 8. Form Handling

```vue
<script setup lang="ts">
const form = reactive({
  name: '',
  email: '',
})

// v-model cho two-way binding
// Validation với computed
const errors = computed(() => ({
  name: !form.name ? 'Name is required' : null,
  email: !form.email.includes('@') ? 'Invalid email' : null,
}))

const isValid = computed(() =>
  Object.values(errors.value).every(e => e === null)
)
</script>

<template>
  <form @submit.prevent="handleSubmit">
    <input v-model="form.name" />
    <span v-if="errors.name" class="error">{{ errors.name }}</span>
  </form>
</template>
```

## 9. Testing

### Component Testing (Vitest + Vue Test Utils)

```typescript
import { mount } from '@vue/test-utils'
import UserCard from './UserCard.vue'

describe('UserCard', () => {
  it('renders user name', () => {
    const wrapper = mount(UserCard, {
      props: { user: { name: 'John', email: 'john@test.com' } },
    })
    expect(wrapper.text()).toContain('John')
  })

  it('emits delete event', async () => {
    const wrapper = mount(UserCard, { props: { user: mockUser } })
    await wrapper.find('[data-testid="delete-btn"]').trigger('click')
    expect(wrapper.emitted('delete')).toBeTruthy()
  })
})
```

### Composable Testing

```typescript
import { useFetch } from './useFetch'

describe('useFetch', () => {
  it('fetches data', async () => {
    const { data, loading } = useFetch('/api/users')
    expect(loading.value).toBe(true)
    await flushPromises()
    expect(data.value).toBeDefined()
    expect(loading.value).toBe(false)
  })
})
```

## 10. Security

- **Không dùng `v-html`** với user input (XSS risk)
- **Sanitize** tất cả user input trước khi render
- **CSP headers** cho production
- **Không store secrets** trong client-side code
- **HTTPS** cho tất cả API calls
