
# Vue Agent


## When to Use

- Building or refactoring Vue 3 components
- Implementing Composition API patterns with `<script setup>`
- Managing state with Pinia
- Structuring composables and component architecture

**Don't use when:** Working in a React, Angular, or Svelte project.

---

## Decision Tree

```
Need reactive primitive?              → ref()
Need reactive object/array?           → reactive() or ref({})
Derived value from state?             → computed()
Watch a ref and run a side effect?    → watch() or watchEffect()
Shared state across components?       → Pinia store
Logic reusable across components?     → composable (useXxx)
Async component or heavy route?       → defineAsyncComponent / lazy route
```

---

## `<script setup>` — Always Use It (REQUIRED)

```vue
<!-- ✅ GOOD — Composition API with <script setup> -->
<script setup lang="ts">
import { ref, computed } from 'vue'

const props = defineProps<{ title: string; count?: number }>()
const emit = defineEmits<{ update: [value: number] }>()

const doubled = computed(() => (props.count ?? 0) * 2)
</script>

<template>
  <h1>{{ title }} — {{ doubled }}</h1>
  <button @click="emit('update', doubled)">Emit</button>
</template>

<!-- ❌ NEVER — Options API in new components -->
<script>
export default {
  props: ['title'],
  data() { return { count: 0 } },
  methods: { increment() { this.count++ } }
}
</script>
```

---

## Critical Patterns

### 1. `ref()` vs `reactive()`

```typescript
// ✅ GOOD — ref() for primitives and single values
const count = ref(0)
const name  = ref('')
const user  = ref<User | null>(null)

count.value++          // access via .value in script
// {{ count }}         // no .value needed in template

// ✅ GOOD — reactive() for objects you always access as a whole
const form = reactive({ name: '', email: '', age: 0 })
form.name = 'Alice'    // no .value needed

// ❌ BAD — reactive() with primitives loses reactivity on reassign
const count = reactive(0) // primitives don't work with reactive()
```

### 2. `computed()` for Derived State

```typescript
// ✅ GOOD — computed caches and tracks dependencies automatically
const fullName  = computed(() => `${firstName.value} ${lastName.value}`)
const activeUsers = computed(() => users.value.filter(u => u.active))

// ✅ GOOD — writable computed
const modelValue = computed({
  get: () => props.value,
  set: (val) => emit('update:modelValue', val),
})

// ❌ BAD — method called on every render, no caching
const fullName = () => `${firstName.value} ${lastName.value}`
```

### 3. `defineProps` and `defineEmits` with TypeScript

```typescript
// ✅ GOOD — type-only declarations (preferred)
const props = defineProps<{
  title: string
  items: Item[]
  loading?: boolean
}>()

const emit = defineEmits<{
  select: [item: Item]
  'update:modelValue': [value: string]
}>()

// ✅ GOOD — with defaults
const props = withDefaults(defineProps<{
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
}>(), {
  size: 'md',
  disabled: false,
})
```

### 4. `v-model` Two-Way Binding

```vue
<!-- ✅ GOOD — v-model in Vue 3 (uses modelValue + update:modelValue) -->
<script setup lang="ts">
defineProps<{ modelValue: string }>()
defineEmits<{ 'update:modelValue': [value: string] }>()
</script>

<template>
  <input
    :value="modelValue"
    @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
  />
</template>

<!-- Usage -->
<MyInput v-model="username" />

<!-- ✅ GOOD — multiple v-model bindings (Vue 3.4+) -->
<UserForm v-model:name="name" v-model:email="email" />
```

### 5. `v-for` Always Needs `:key`

```vue
<!-- ✅ GOOD — stable, unique key -->
<li v-for="item in items" :key="item.id">{{ item.name }}</li>

<!-- ✅ GOOD — when items have no id, use index only if list is static -->
<li v-for="(item, index) in staticList" :key="index">{{ item }}</li>

<!-- ❌ BAD — no key, Vue can't optimize re-renders -->
<li v-for="item in items">{{ item.name }}</li>

<!-- ❌ BAD — index key on dynamic list causes incorrect reuse -->
<li v-for="(item, index) in dynamicList" :key="index">{{ item.name }}</li>
```

### 6. Composables — Reusable Logic

```typescript
// composables/useUser.ts
// ✅ GOOD — composable encapsulates reactive logic
import { ref, onMounted } from 'vue'
import type { User } from '@/types'

export function useUser(userId: string) {
  const user    = ref<User | null>(null)
  const loading = ref(false)
  const error   = ref<Error | null>(null)

  async function fetch() {
    loading.value = true
    error.value   = null
    try {
      user.value = await api.getUser(userId)
    } catch (e) {
      error.value = e as Error
    } finally {
      loading.value = false
    }
  }

  onMounted(fetch)

  return { user, loading, error, refresh: fetch }
}

// Usage in component
const { user, loading, error } = useUser(props.userId)
```

### 7. Pinia — State Management

```typescript
// stores/user.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

// ✅ GOOD — Setup Store syntax (Composition API style)
export const useUserStore = defineStore('user', () => {
  const currentUser = ref<User | null>(null)
  const isAuthenticated = computed(() => currentUser.value !== null)

  async function login(credentials: Credentials) {
    currentUser.value = await authService.login(credentials)
  }

  function logout() {
    currentUser.value = null
  }

  return { currentUser, isAuthenticated, login, logout }
})

// Usage in component
const userStore = useUserStore()
const { currentUser, isAuthenticated } = storeToRefs(userStore) // keep reactivity!
```

```typescript
// ❌ BAD — destructuring store directly loses reactivity
const { currentUser } = useUserStore() // currentUser is no longer reactive
```

### 8. Template Refs

```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'

// ✅ GOOD — typed template ref
const inputEl = ref<HTMLInputElement | null>(null)

onMounted(() => {
  inputEl.value?.focus()
})
</script>

<template>
  <input ref="inputEl" type="text" />
</template>
```

### 9. `watch` vs `watchEffect`

```typescript
// ✅ watchEffect — runs immediately, tracks deps automatically
watchEffect(() => {
  document.title = `${route.name as string} — MyApp`
})

// ✅ watch — explicit source, runs only on change, has old value
watch(userId, async (newId, oldId) => {
  if (newId !== oldId) await fetchUser(newId)
})

// ✅ watch with immediate
watch(userId, fetchUser, { immediate: true })

// ❌ BAD — watch on reactive object without deep or specific getter
watch(form, handler)             // may not trigger on nested changes
watch(() => form, handler)       // always returns same reference

// ✅ GOOD — watch specific property of reactive
watch(() => form.email, validateEmail)
```

### 10. Async Components & Lazy Routes

```typescript
// ✅ GOOD — lazy load heavy components
import { defineAsyncComponent } from 'vue'

const HeavyChart = defineAsyncComponent(() =>
  import('@/components/HeavyChart.vue')
)

// ✅ GOOD — lazy routes in Vue Router 4
const routes = [
  {
    path: '/dashboard',
    component: () => import('@/views/DashboardView.vue'),
  },
]
```

---

## Component Structure Convention

```
src/
├── components/          # Reusable presentational components
│   └── BaseButton/
│       ├── BaseButton.vue
│       └── index.ts
├── composables/         # Shared stateful logic (useXxx.ts)
├── stores/              # Pinia stores (useXxxStore.ts)
├── views/               # Route-level components (pages)
├── router/              # Vue Router config
└── types/               # TypeScript interfaces and types
```

**Rules:**
- Components in `components/` are prefix-named: `BaseButton`, `AppHeader`, `TheNavbar`.
- `The` prefix = singleton (one instance per app), `Base` prefix = highly reusable.
- Composables always start with `use`: `useAuth`, `useCart`.
- Stores always end with `Store`: `useUserStore`, `useCartStore`.

---

## Commands

```bash
npm create vue@latest          # scaffold new project (Vite + TS + Router + Pinia)
npx vue-tsc --noEmit           # type-check templates
npm run lint                   # ESLint with vue plugin
npm run build                  # production build
```

---

## Definition of Done

- [ ] All components use `<script setup lang="ts">`
- [ ] No Options API in new code
- [ ] All `v-for` loops have a stable `:key`
- [ ] Pinia store properties destructured with `storeToRefs()`
- [ ] Composables named `useXxx` and return reactive refs
- [ ] `computed()` used for all derived state (not methods)
- [ ] No direct DOM manipulation — use template refs + lifecycle hooks
- [ ] Async components or lazy routes used for heavy views
