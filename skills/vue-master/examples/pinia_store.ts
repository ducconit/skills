/**
 * Pinia Store Example — Setup Store Syntax
 *
 * Demonstrates:
 * - Setup store (Composition API style)
 * - Typed state, getters, actions
 * - Async actions with error handling
 * - Store composition (using one store in another)
 */

import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

// ─── Types ──────────────────────────────────────────────────────

interface User {
  id: string
  name: string
  email: string
  role: 'admin' | 'user' | 'guest'
}

interface LoginCredentials {
  email: string
  password: string
}

interface AuthTokens {
  accessToken: string
  refreshToken: string
  expiresAt: number
}

// ─── Auth Store ─────────────────────────────────────────────────

export const useAuthStore = defineStore('auth', () => {
  // ── State ──
  const user = ref<User | null>(null)
  const tokens = ref<AuthTokens | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // ── Getters (computed) ──
  const isAuthenticated = computed(() => !!tokens.value?.accessToken)

  const isAdmin = computed(() => user.value?.role === 'admin')

  const displayName = computed(() => user.value?.name ?? 'Guest')

  const isTokenExpired = computed(() => {
    if (!tokens.value) return true
    return Date.now() > tokens.value.expiresAt
  })

  // ── Actions ──
  async function login(credentials: LoginCredentials): Promise<void> {
    loading.value = true
    error.value = null

    try {
      // API call (replace with actual implementation)
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(credentials),
      })

      if (!response.ok) {
        throw new Error('Invalid credentials')
      }

      const data = await response.json()
      user.value = data.user
      tokens.value = data.tokens
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Login failed'
      throw e
    } finally {
      loading.value = false
    }
  }

  function logout(): void {
    user.value = null
    tokens.value = null
    error.value = null
  }

  async function refreshToken(): Promise<void> {
    if (!tokens.value?.refreshToken) {
      logout()
      return
    }

    try {
      const response = await fetch('/api/auth/refresh', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refreshToken: tokens.value.refreshToken }),
      })

      if (!response.ok) {
        throw new Error('Token refresh failed')
      }

      const data = await response.json()
      tokens.value = data.tokens
    } catch {
      logout()
    }
  }

  return {
    // State
    user,
    tokens,
    loading,
    error,
    // Getters
    isAuthenticated,
    isAdmin,
    displayName,
    isTokenExpired,
    // Actions
    login,
    logout,
    refreshToken,
  }
})

// ─── Notification Store ─────────────────────────────────────────

interface Notification {
  id: string
  message: string
  type: 'success' | 'error' | 'info' | 'warning'
  createdAt: number
}

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref<Notification[]>([])

  const unreadCount = computed(() => notifications.value.length)

  function add(message: string, type: Notification['type'] = 'info') {
    const notification: Notification = {
      id: crypto.randomUUID(),
      message,
      type,
      createdAt: Date.now(),
    }
    notifications.value.push(notification)

    // Auto-remove after 5 seconds
    setTimeout(() => remove(notification.id), 5000)
  }

  function remove(id: string) {
    notifications.value = notifications.value.filter((n) => n.id !== id)
  }

  function clear() {
    notifications.value = []
  }

  return { notifications, unreadCount, add, remove, clear }
})

// ─── Cart Store (Store Composition Example) ─────────────────────

interface CartItem {
  productId: string
  name: string
  price: number
  quantity: number
}

export const useCartStore = defineStore('cart', () => {
  // ── Using another store ──
  const authStore = useAuthStore()
  const notificationStore = useNotificationStore()

  // ── State ──
  const items = ref<CartItem[]>([])

  // ── Getters ──
  const totalItems = computed(() =>
    items.value.reduce((sum, item) => sum + item.quantity, 0)
  )

  const totalPrice = computed(() =>
    items.value.reduce((sum, item) => sum + item.price * item.quantity, 0)
  )

  const isEmpty = computed(() => items.value.length === 0)

  // ── Actions ──
  function addItem(item: Omit<CartItem, 'quantity'>) {
    const existing = items.value.find((i) => i.productId === item.productId)
    if (existing) {
      existing.quantity++
    } else {
      items.value.push({ ...item, quantity: 1 })
    }
    notificationStore.add(`${item.name} added to cart`, 'success')
  }

  function removeItem(productId: string) {
    items.value = items.value.filter((i) => i.productId !== productId)
  }

  function updateQuantity(productId: string, quantity: number) {
    const item = items.value.find((i) => i.productId === productId)
    if (item) {
      if (quantity <= 0) {
        removeItem(productId)
      } else {
        item.quantity = quantity
      }
    }
  }

  async function checkout() {
    // Require authentication
    if (!authStore.isAuthenticated) {
      notificationStore.add('Please login to checkout', 'warning')
      throw new Error('Authentication required')
    }

    try {
      await fetch('/api/orders', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${authStore.tokens?.accessToken}`,
        },
        body: JSON.stringify({ items: items.value }),
      })

      items.value = []
      notificationStore.add('Order placed successfully!', 'success')
    } catch (e) {
      notificationStore.add('Checkout failed', 'error')
      throw e
    }
  }

  return {
    items,
    totalItems,
    totalPrice,
    isEmpty,
    addItem,
    removeItem,
    updateQuantity,
    checkout,
  }
})
