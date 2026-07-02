/**
 * Composable Pattern Example — useFetch
 *
 * Demonstrates:
 * - Creating a reusable composable
 * - Proper TypeScript typing with generics
 * - Cleanup with onUnmounted (AbortController)
 * - Reactive return values
 * - MaybeRefOrGetter for flexible input
 */

import {
  ref,
  watchEffect,
  toValue,
  onUnmounted,
  type Ref,
  type MaybeRefOrGetter,
} from 'vue'

// ─── Types ──────────────────────────────────────────────────────

interface UseFetchOptions {
  /** Fetch immediately on mount */
  immediate?: boolean
  /** Custom fetch init options */
  fetchOptions?: RequestInit
  /** Transform response before setting data */
  transform?: <T>(data: unknown) => T
}

interface UseFetchReturn<T> {
  data: Ref<T | null>
  error: Ref<Error | null>
  loading: Ref<boolean>
  /** Re-execute the fetch */
  refresh: () => Promise<void>
  /** Abort the current request */
  abort: () => void
}

// ─── Composable ─────────────────────────────────────────────────

export function useFetch<T = unknown>(
  url: MaybeRefOrGetter<string>,
  options: UseFetchOptions = {}
): UseFetchReturn<T> {
  const { immediate = true, fetchOptions = {}, transform } = options

  // Reactive state
  const data = ref<T | null>(null) as Ref<T | null>
  const error = ref<Error | null>(null)
  const loading = ref(false)

  // AbortController for cleanup
  let abortController: AbortController | null = null

  function abort() {
    abortController?.abort()
    abortController = null
  }

  async function execute() {
    // Abort previous request
    abort()

    // Create new AbortController
    abortController = new AbortController()

    loading.value = true
    error.value = null

    try {
      const response = await fetch(toValue(url), {
        ...fetchOptions,
        signal: abortController.signal,
      })

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`)
      }

      const json = await response.json()
      data.value = transform ? transform<T>(json) : (json as T)
    } catch (e) {
      // Ignore abort errors
      if (e instanceof DOMException && e.name === 'AbortError') {
        return
      }
      error.value = e instanceof Error ? e : new Error(String(e))
    } finally {
      loading.value = false
    }
  }

  // Auto-fetch when URL changes (if reactive)
  if (immediate) {
    watchEffect(() => {
      // toValue(url) creates a dependency on the URL
      toValue(url)
      execute()
    })
  }

  // Cleanup on unmount
  onUnmounted(() => {
    abort()
  })

  return {
    data,
    error,
    loading,
    refresh: execute,
    abort,
  }
}

// ─── Usage Examples ─────────────────────────────────────────────

/*
// Basic usage
const { data: users, loading, error } = useFetch<User[]>('/api/users')

// With reactive URL
const userId = ref('123')
const { data: user } = useFetch<User>(() => `/api/users/${userId.value}`)

// With options
const { data, refresh } = useFetch<Product[]>('/api/products', {
  immediate: false,
  fetchOptions: {
    headers: { 'Authorization': 'Bearer token' },
  },
  transform: (raw) => (raw as any).data,
})

// Manual fetch
await refresh()
*/
