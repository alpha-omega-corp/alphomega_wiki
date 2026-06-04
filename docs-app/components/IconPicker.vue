<script setup lang="ts">
/**
 * Lucide icon picker. Emits a Nuxt UI icon name (e.g. `i-lucide-flame`).
 *
 * Fully offline: the name list comes from our own `/api/icons/lucide` route,
 * which reads the locally-installed `@iconify-json/lucide` collection (no
 * external network), and the icons themselves render from Nuxt Icon's local
 * server bundle. The list is cached in `useState` so re-opening is instant.
 */
const props = defineProps<{ modelValue?: string }>()
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

// Open state is a model so callers can also open the picker programmatically
// (e.g. from a context-menu action), not only by clicking the trigger.
const open = defineModel<boolean>('open', { default: false })
const query = ref('')

// Curated fallback used only if the Iconify collection can't be fetched.
const FALLBACK = [
  'info', 'circle-help', 'triangle-alert', 'octagon-alert', 'check', 'check-circle',
  'x', 'flame', 'star', 'heart', 'bookmark', 'flag', 'bell', 'zap', 'lightbulb',
  'rocket', 'sparkles', 'book-open', 'file-text', 'code', 'terminal', 'settings',
  'wrench', 'shield', 'lock', 'key', 'cloud', 'database', 'server', 'cpu', 'bug',
]

// Shared across all picker instances + SSR-safe.
const names = useState<string[]>('lucide-icon-names', () => [])
const pending = ref(false)
const failed = ref(false)

async function load() {
  if (names.value.length || pending.value) return
  pending.value = true
  failed.value = false
  try {
    const data = await $fetch<string[]>('/api/icons/lucide')
    names.value = data?.length ? data : FALLBACK
  } catch {
    failed.value = true
    names.value = FALLBACK
  } finally {
    pending.value = false
  }
}

// Load on the client; the local route reads the bundled collection (offline).
onMounted(load)

const filtered = computed(() => {
  const q = query.value.trim().toLowerCase()
  return q ? names.value.filter(n => n.includes(q)) : names.value
})

// Cap how many we render at once to keep the grid snappy.
const LIMIT = 96
const visible = computed(() => filtered.value.slice(0, LIMIT))

function select(name: string) {
  emit('update:modelValue', `i-lucide-${name}`)
  open.value = false
}

function clear() {
  emit('update:modelValue', '')
  open.value = false
}
</script>

<template>
  <UModal v-model:open="open" title="Pick an icon" :ui="{ content: 'sm:max-w-lg' }">
    <!-- Trigger is overridable so callers can use e.g. the card's icon itself. -->
    <slot name="trigger" :icon="modelValue">
      <UButton variant="outline" color="neutral" size="xs" class="gap-1">
        <UIcon v-if="modelValue" :name="modelValue" class="size-4"/>
        <UIcon v-else name="i-lucide-smile" class="size-4 text-muted"/>
        <UIcon name="i-lucide-chevron-down" class="size-3 text-muted"/>
      </UButton>
    </slot>

    <template #body>
      <div class="flex items-center gap-1 mb-3">
        <UInput
            v-model="query"
            icon="i-lucide-search"
            placeholder="Search Lucide icons"
            autofocus
            class="flex-1"
        />
        <UButton
            v-if="modelValue"
            icon="i-lucide-eraser"
            variant="ghost"
            color="neutral"
            aria-label="Clear icon"
            @click="clear"
        />
      </div>

      <div v-if="pending" class="py-10 text-center text-sm text-muted">Loading icons…</div>

      <template v-else>
        <p v-if="failed" class="mb-2 text-xs text-warning">Offline — showing a limited set.</p>

        <div v-if="visible.length" class="grid grid-cols-10 gap-0.5 max-h-[55vh] overflow-y-auto">
          <UButton
              v-for="name in visible"
              :key="name"
              :icon="`i-lucide-${name}`"
              variant="ghost"
              color="neutral"
              square
              :title="name"
              :class="modelValue === `i-lucide-${name}` ? 'bg-primary-50 dark:bg-primary-900/30' : ''"
              @click="select(name)"
          />
        </div>
        <div v-else class="py-10 text-center text-sm text-muted">No icons match “{{ query }}”.</div>

        <p v-if="filtered.length > visible.length" class="mt-3 text-xs text-muted text-center">
          Showing {{ visible.length }} of {{ filtered.length }} — refine your search.
        </p>
      </template>
    </template>
  </UModal>
</template>
