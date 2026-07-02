<!--
  Typed Component Example — UserCard

  Demonstrates:
  - <script setup lang="ts"> with full TypeScript
  - defineProps with interface
  - defineEmits with typed events
  - Computed properties
  - Template refs typed
  - CSS scoped
-->

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

// ─── Types ──────────────────────────────────────────────────────

interface User {
  id: string
  name: string
  email: string
  avatar?: string
  role: 'admin' | 'user' | 'guest'
  createdAt: string
}

// ─── Props ──────────────────────────────────────────────────────

interface Props {
  user: User
  showActions?: boolean
  highlighted?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showActions: true,
  highlighted: false,
})

// ─── Emits ──────────────────────────────────────────────────────

const emit = defineEmits<{
  edit: [user: User]
  delete: [userId: string]
  'update:highlighted': [value: boolean]
}>()

// ─── Template Refs ──────────────────────────────────────────────

const cardRef = ref<HTMLDivElement | null>(null)
const nameInputRef = ref<HTMLInputElement | null>(null)

// ─── Reactive State ─────────────────────────────────────────────

const isEditing = ref(false)
const editName = ref('')

// ─── Computed ───────────────────────────────────────────────────

const initials = computed(() => {
  return props.user.name
    .split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
})

const roleBadgeClass = computed(() => {
  const classMap: Record<User['role'], string> = {
    admin: 'badge--admin',
    user: 'badge--user',
    guest: 'badge--guest',
  }
  return classMap[props.user.role]
})

const memberSince = computed(() => {
  return new Date(props.user.createdAt).toLocaleDateString('vi-VN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
})

// ─── Methods ────────────────────────────────────────────────────

function onEdit() {
  isEditing.value = true
  editName.value = props.user.name
  // Focus input after DOM update
  setTimeout(() => nameInputRef.value?.focus(), 0)
}

function onSave() {
  isEditing.value = false
  emit('edit', { ...props.user, name: editName.value })
}

function onCancel() {
  isEditing.value = false
  editName.value = ''
}

function onDelete() {
  emit('delete', props.user.id)
}

// ─── Lifecycle ──────────────────────────────────────────────────

onMounted(() => {
  // Scroll card into view if highlighted
  if (props.highlighted && cardRef.value) {
    cardRef.value.scrollIntoView({ behavior: 'smooth', block: 'center' })
  }
})
</script>

<template>
  <div
    ref="cardRef"
    class="user-card"
    :class="{ 'user-card--highlighted': highlighted }"
    data-testid="user-card"
  >
    <!-- Avatar -->
    <div class="user-card__avatar">
      <img
        v-if="user.avatar"
        :src="user.avatar"
        :alt="`Avatar of ${user.name}`"
        class="user-card__avatar-img"
      />
      <span v-else class="user-card__avatar-initials">
        {{ initials }}
      </span>
    </div>

    <!-- Info -->
    <div class="user-card__info">
      <div class="user-card__name">
        <template v-if="isEditing">
          <input
            ref="nameInputRef"
            v-model="editName"
            class="user-card__name-input"
            @keyup.enter="onSave"
            @keyup.escape="onCancel"
          />
        </template>
        <template v-else>
          <h3>{{ user.name }}</h3>
        </template>
        <span class="badge" :class="roleBadgeClass">
          {{ user.role }}
        </span>
      </div>
      <p class="user-card__email">{{ user.email }}</p>
      <p class="user-card__date">Member since {{ memberSince }}</p>
    </div>

    <!-- Actions -->
    <div v-if="showActions" class="user-card__actions">
      <template v-if="isEditing">
        <button data-testid="save-btn" @click="onSave">Save</button>
        <button data-testid="cancel-btn" @click="onCancel">Cancel</button>
      </template>
      <template v-else>
        <button data-testid="edit-btn" @click="onEdit">Edit</button>
        <button data-testid="delete-btn" @click="onDelete">Delete</button>
      </template>
    </div>
  </div>
</template>

<style scoped>
.user-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  border: 1px solid var(--border-color, #e2e8f0);
  border-radius: 0.5rem;
  transition: box-shadow 0.2s ease;
}

.user-card:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.user-card--highlighted {
  border-color: var(--primary-color, #3b82f6);
  background-color: var(--primary-bg, #eff6ff);
}

.user-card__avatar {
  flex-shrink: 0;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  overflow: hidden;
  background-color: var(--primary-color, #3b82f6);
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-card__avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.user-card__avatar-initials {
  color: white;
  font-weight: 600;
  font-size: 0.875rem;
}

.user-card__info {
  flex: 1;
  min-width: 0;
}

.user-card__name {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.user-card__name h3 {
  margin: 0;
  font-size: 1rem;
}

.user-card__email {
  margin: 0.25rem 0 0;
  color: var(--text-secondary, #64748b);
  font-size: 0.875rem;
}

.user-card__date {
  margin: 0.125rem 0 0;
  color: var(--text-tertiary, #94a3b8);
  font-size: 0.75rem;
}

.badge {
  padding: 0.125rem 0.5rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 500;
}

.badge--admin {
  background-color: #fee2e2;
  color: #dc2626;
}

.badge--user {
  background-color: #dbeafe;
  color: #2563eb;
}

.badge--guest {
  background-color: #f1f5f9;
  color: #64748b;
}

.user-card__actions {
  display: flex;
  gap: 0.5rem;
}

.user-card__name-input {
  padding: 0.25rem 0.5rem;
  border: 1px solid var(--border-color, #e2e8f0);
  border-radius: 0.25rem;
  font-size: 1rem;
}
</style>
