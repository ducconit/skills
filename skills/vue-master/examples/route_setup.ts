/**
 * Vue Router Setup Example
 *
 * Demonstrates:
 * - Lazy-loaded routes with code splitting
 * - Navigation guards (global, per-route)
 * - Meta fields with TypeScript
 * - Nested routes
 * - Route-level middleware
 */

import {
  createRouter,
  createWebHistory,
  type RouteRecordRaw,
  type RouteLocationNormalized,
} from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// ─── Route Meta Types ───────────────────────────────────────────

// Extend RouteMeta for type-safe meta fields
declare module 'vue-router' {
  interface RouteMeta {
    /** Require authentication */
    requiresAuth?: boolean
    /** Required user roles */
    roles?: Array<'admin' | 'user' | 'guest'>
    /** Page title for document.title */
    title?: string
    /** Layout to use */
    layout?: 'default' | 'auth' | 'blank'
    /** Transition name */
    transition?: string
  }
}

// ─── Routes ─────────────────────────────────────────────────────

const routes: RouteRecordRaw[] = [
  // Public routes
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/HomePage.vue'),
    meta: {
      title: 'Home',
      layout: 'default',
    },
  },

  // Auth routes
  {
    path: '/auth',
    meta: { layout: 'auth' },
    children: [
      {
        path: 'login',
        name: 'login',
        component: () => import('@/features/auth/views/LoginPage.vue'),
        meta: { title: 'Login' },
      },
      {
        path: 'register',
        name: 'register',
        component: () => import('@/features/auth/views/RegisterPage.vue'),
        meta: { title: 'Register' },
      },
    ],
  },

  // Protected routes with nested layout
  {
    path: '/dashboard',
    component: () => import('@/features/dashboard/layouts/DashboardLayout.vue'),
    meta: {
      requiresAuth: true,
      layout: 'default',
    },
    children: [
      {
        path: '',
        name: 'dashboard',
        component: () => import('@/features/dashboard/views/DashboardHome.vue'),
        meta: { title: 'Dashboard' },
      },
      {
        path: 'analytics',
        name: 'analytics',
        component: () => import('@/features/dashboard/views/AnalyticsPage.vue'),
        meta: {
          title: 'Analytics',
          roles: ['admin'],
        },
      },
      {
        path: 'settings',
        name: 'settings',
        component: () => import('@/features/dashboard/views/SettingsPage.vue'),
        meta: { title: 'Settings' },
      },
    ],
  },

  // Dynamic route with params
  {
    path: '/users/:id',
    name: 'user-profile',
    component: () => import('@/features/users/views/UserProfile.vue'),
    meta: {
      requiresAuth: true,
      title: 'User Profile',
    },
    // Per-route guard
    beforeEnter: (to) => {
      const userId = to.params.id as string
      if (!userId || userId === 'undefined') {
        return { name: 'home' }
      }
    },
  },

  // 404 catch-all (must be last)
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/views/NotFoundPage.vue'),
    meta: {
      title: '404 Not Found',
      layout: 'blank',
    },
  },
]

// ─── Router Instance ────────────────────────────────────────────

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  // Scroll behavior
  scrollBehavior(to, _from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    if (to.hash) {
      return { el: to.hash, behavior: 'smooth' }
    }
    return { top: 0 }
  },
})

// ─── Global Navigation Guards ───────────────────────────────────

// Authentication guard
router.beforeEach(async (to: RouteLocationNormalized) => {
  const authStore = useAuthStore()

  // Set document title
  const title = to.meta.title
  document.title = title ? `${title} | MyApp` : 'MyApp'

  // Check authentication
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return {
      name: 'login',
      query: { redirect: to.fullPath },
    }
  }

  // Check token expiry
  if (authStore.isAuthenticated && authStore.isTokenExpired) {
    try {
      await authStore.refreshToken()
    } catch {
      return { name: 'login' }
    }
  }

  // Check role-based access
  if (to.meta.roles && to.meta.roles.length > 0) {
    const userRole = authStore.user?.role
    if (!userRole || !to.meta.roles.includes(userRole)) {
      return { name: 'dashboard' } // Redirect to dashboard if unauthorized
    }
  }

  // Redirect authenticated users away from auth pages
  if (to.path.startsWith('/auth') && authStore.isAuthenticated) {
    return { name: 'dashboard' }
  }
})

// Loading state management
router.beforeEach(() => {
  // Start progress bar (e.g., NProgress)
  // NProgress.start()
})

router.afterEach(() => {
  // End progress bar
  // NProgress.done()
})

// Error handling
router.onError((error) => {
  // Handle chunk loading errors (e.g., after deployment)
  if (
    error.message.includes('Failed to fetch dynamically imported module') ||
    error.message.includes('Loading chunk')
  ) {
    window.location.reload()
  }
})

export default router
