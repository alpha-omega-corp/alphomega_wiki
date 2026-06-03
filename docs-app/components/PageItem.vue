<template>
  <li>
    <UContextMenu :items="menuItems">
      <div class="flex items-center gap-0.5 group">
        <UButton
            v-if="hasChildren"
            :icon="open ? 'i-lucide-chevron-down' : 'i-lucide-chevron-right'"
            size="xs"
            variant="ghost"
            color="neutral"
            class="shrink-0"
            @click="open = !open"
        />
        <span v-else class="size-6 shrink-0"/>

        <UInput
            v-if="renaming"
            v-model="draftTitle"
            size="xs"
            autofocus
            class="flex-1"
            @keydown.enter.prevent="confirmRename"
            @keydown.esc="cancelRename"
            @blur="confirmRename"
        />
        <NuxtLink
            v-else
            :to="`/docs/${fullSlug}`"
            class="flex flex-1 items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors min-w-0"
            :class="isActive
            ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400 font-medium'
            : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-800'"
        >
          <span class="truncate">{{ page.title }}</span>
        </NuxtLink>

        <UButton
            icon="i-lucide-plus"
            size="xs"
            variant="ghost"
            color="neutral"
            class="shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
            @click="$emit('create', page.id)"
        />

        <UButton
            icon="i-lucide-more-horizontal"
            size="xs"
            variant="ghost"
            color="neutral"
            class="shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
            aria-label="Page options"
            @click="openMenu"
        />
      </div>
    </UContextMenu>

    <ul
        v-if="hasChildren"
        v-show="open"
        class="mt-0.5 ml-3 space-y-0.5 border-l border-gray-200 dark:border-gray-700 pl-2"
    >
      <PageItem
          v-for="child in page.pages"
          :key="child.id"
          :page="child"
          :base-path="fullSlug"
          @create="$emit('create', $event)"
      />
    </ul>
  </li>
</template>

<script setup lang="ts">
import type {Page} from "~/types/page";
import type {ContextMenuItem} from "@nuxt/ui";

const props = withDefaults(defineProps<{ page: Page; basePath?: string }>(), {
  basePath: '',
})
const emit = defineEmits<{ create: [id: number] }>()

// Full slug path including ancestors, e.g. `parent/child`.
const fullSlug = computed(() =>
    props.basePath ? `${props.basePath}/${props.page.slug}` : props.page.slug,
)

const route = useRoute()
const toast = useToast()

// Context-menu actions for this page (rendered as ContextMenuItems).
const menuItems = computed<ContextMenuItem[][]>(() => [
  [
    {
      label: 'Create Page',
      icon: 'i-lucide-file-plus',
      onSelect: () => emit('create', props.page.id),
    },
    {
      label: 'Rename',
      icon: 'i-lucide-pencil',
      onSelect: () => startRename(),
    },
    {
      label: 'Export',
      icon: 'i-lucide-download',
      onSelect: () => exportMarkdown(),
    },
  ],
  [
    {
      label: 'Delete',
      icon: 'i-lucide-trash-2',
      color: 'error',
      onSelect: () => deletePage(),
    },
  ],
])

// UContextMenu opens on right-click; the three-dots button is a left click, so
// synthesize a contextmenu event on it to open the menu at the button.
function openMenu(event: MouseEvent) {
  const el = event.currentTarget as HTMLElement
  const rect = el.getBoundingClientRect()
  el.dispatchEvent(new MouseEvent('contextmenu', {
    bubbles: true,
    clientX: rect.left,
    clientY: rect.bottom,
  }))
}

// Export the page's markdown as a downloadable .md file.
async function exportMarkdown() {
  try {
    // Sidebar list items don't carry the body, so fetch the full page to get
    // the markdown content at export time.
    const {content} = await $fetch<Page>(`/api/page/${props.page.id}`)

    const blob = new Blob([content ?? ''], {type: 'text/markdown;charset=utf-8'})
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${props.page.slug}.md`
    link.click()
    URL.revokeObjectURL(url)
  } catch (e: any) {
    toast.add({
      title: 'Export failed',
      description: e?.data?.message ?? e?.message ?? 'Could not export the page.',
      color: 'error',
      icon: 'i-lucide-triangle-alert',
    })
  }
}

// Delete the page (after confirmation), then refresh the sidebar tree.
async function deletePage() {
  if (!confirm(`Delete "${props.page.title}"? This cannot be undone.`)) return

  try {
    await $fetch(`/api/page/${props.page.id}`, {method: 'DELETE'})
    await refreshNuxtData('pages')
    toast.add({
      title: 'Page deleted',
      description: props.page.title,
      color: 'success',
      icon: 'i-lucide-trash-2',
    })
    // If we just deleted the page being viewed, leave the now-dead route.
    if (isActive.value) await navigateTo('/docs/getting-started')
  } catch (e: any) {
    toast.add({
      title: 'Delete failed',
      description: e?.data?.message ?? e?.message ?? 'Could not delete the page.',
      color: 'error',
      icon: 'i-lucide-triangle-alert',
    })
  }
}

// Inline rename.
const renaming = ref(false)
const draftTitle = ref('')

function startRename() {
  draftTitle.value = props.page.title
  renaming.value = true
}

function cancelRename() {
  renaming.value = false
}

async function confirmRename() {
  if (!renaming.value) return
  const title = draftTitle.value.trim()
  renaming.value = false
  if (!title || title === props.page.title) return

  try {
    await $fetch(`/api/page/${props.page.id}`, {
      method: 'PUT',
      body: {title},
    })
    await refreshNuxtData('pages').then(() => {
      toast.add({
        title: 'Page Renamed',
        description: `${props.page.title}`,
        color: 'success',
        icon: 'i-lucide-checkmark',
      })
    })
  } catch (e: any) {
    toast.add({
      title: 'Rename failed',
      description: e?.data?.message ?? e?.message ?? 'Could not rename the page.',
      color: 'error',
      icon: 'i-lucide-triangle-alert',
    })
  }
}

// Expanded state is shared app-wide and keyed by page id so it survives
// navigation between docs pages (and any component remount).
const expanded = useState<Record<number, boolean>>('pagelist-expanded', () => ({}))
const open = computed({
  get: () => expanded.value[props.page.id] ?? false,
  set: (value) => {
    expanded.value[props.page.id] = value
  },
})

const hasChildren = computed(() => !!props.page.pages?.length)
const isActive = computed(() => route.path === `/docs/${fullSlug.value}`)

function pageInSubtree(page: Page, base: string, path: string): boolean {
  const full = base ? `${base}/${page.slug}` : page.slug
  if (`/docs/${full}` === path) return true
  return (page.pages ?? []).some(p => pageInSubtree(p, full, path))
}

const isActiveSubtree = computed(() => pageInSubtree(props.page, props.basePath, route.path))

watchEffect(() => {
  if (isActiveSubtree.value) open.value = true
})
</script>