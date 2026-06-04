<template>
  <li>
    <UContextMenu :items="menuItems">
      <div
          ref="rowRef"
          class="relative flex items-center gap-0.5 group rounded-md transition-colors"
          :class="{
            'opacity-50': isDragging,
            'bg-primary-50 dark:bg-primary-900/20 ring-1 ring-primary-400 ring-inset': dropPos === 'inside',
          }"
          :draggable="!renaming && loggedIn"
          @dragstart.stop="onDragStart"
          @dragend="onDragEnd"
          @dragover.prevent="onDragOver"
          @dragleave="onDragLeave"
          @drop.prevent="onDrop"
      >
        <!-- Sibling drop indicators (a line above / below the row). -->
        <div v-if="dropPos === 'before'" class="absolute -top-px left-6 right-1 h-0.5 rounded-full bg-primary-500"/>
        <div v-if="dropPos === 'after'" class="absolute -bottom-px left-6 right-1 h-0.5 rounded-full bg-primary-500"/>

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

        <!-- Page icon doubles as the icon-picker trigger; also opened via the context menu. -->
        <IconPicker
            v-model:open="iconPickerOpen"
            :model-value="page.icon ?? ''"
            @update:model-value="onIconSelect"
        >
          <template #trigger>
            <span
                v-if="page.icon"
                class="flex size-5 shrink-0 cursor-pointer items-center justify-center"
                title="Change icon"
                @click.stop.prevent
            >
              <UIcon v-if="isIcon(page.icon)" :name="page.icon" class="size-4 text-gray-500 dark:text-gray-400"/>
              <span v-else class="text-sm leading-none">{{ page.icon }}</span>
            </span>
            <span v-else class="block w-0" aria-hidden="true"/>
          </template>
        </IconPicker>

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
            draggable="false"
            class="flex flex-1 items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors min-w-0"
            :class="isActive
            ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400 font-medium'
            : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-800'"
        >
          <span class="truncate">{{ page.title }}</span>
        </NuxtLink>

        <UButton
            v-if="loggedIn"
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
import {PageDndKey, type DropPosition} from "~/composables/pageDnd";

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

// Editing actions are gated behind login (the server enforces this too).
const {loggedIn} = useAuth()

// Context-menu actions for this page (rendered as ContextMenuItems). Export is
// a read-only, client-side action so it stays available when logged out.
const menuItems = computed<ContextMenuItem[][]>(() => {
  const exportItem: ContextMenuItem = {
    label: 'Export',
    icon: 'i-lucide-download',
    onSelect: () => exportMarkdown(),
  }
  if (!loggedIn.value) return [[exportItem]]
  return [
    [
      {label: 'Create Page', icon: 'i-lucide-file-plus', onSelect: () => emit('create', props.page.id)},
      {label: 'Rename', icon: 'i-lucide-pencil', onSelect: () => startRename()},
      {label: 'Change Icon', icon: 'i-lucide-image', onSelect: () => openIconPicker()},
      exportItem,
    ],
    [
      {label: 'Delete', icon: 'i-lucide-trash-2', color: 'error', onSelect: () => deletePage()},
    ],
  ]
})

// --- Icon -----------------------------------------------------------------
const iconPickerOpen = ref(false)

function isIcon(value: string) {
  return /^i-[a-z0-9]+[:-]/.test(value) || value.includes(':')
}

// Open the icon-selection modal (deferred a tick so it doesn't clash with the
// closing context menu's focus handling).
async function openIconPicker() {
  await nextTick()
  iconPickerOpen.value = true
}

// Persist the chosen icon (empty string clears it). Updates the tree node
// optimistically so it renders immediately, then re-syncs with the server.
async function onIconSelect(icon: string) {
  const previous = props.page.icon
  props.page.icon = icon
  try {
    await $fetch(`/api/page/${props.page.id}`, {method: 'PUT', body: {icon}})
    await refreshNuxtData('pages')
  } catch (e: any) {
    props.page.icon = previous
    toast.add({
      title: 'Could not update icon',
      description: e?.data?.message ?? e?.message ?? 'Failed to set the page icon.',
      color: 'error',
      icon: 'i-lucide-triangle-alert',
    })
  }
}

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

// --- Drag & drop -----------------------------------------------------------
const dnd = inject(PageDndKey)!
const rowRef = ref<HTMLElement>()
// Where a drop would land on *this* row, used to render the indicator.
const dropPos = ref<DropPosition | null>(null)
const isDragging = computed(() => dnd.draggingId.value === props.page.id)

function onDragStart(event: DragEvent) {
  if (event.dataTransfer) {
    event.dataTransfer.setData('text/plain', String(props.page.id))
    event.dataTransfer.effectAllowed = 'move'
  }
  dnd.start(props.page.id)
}

function onDragEnd() {
  dnd.end()
  dropPos.value = null
}

// Top quarter → drop before (sibling), bottom quarter → after (sibling),
// middle → inside (nest as a child).
function onDragOver(event: DragEvent) {
  const id = dnd.draggingId.value
  if (id === null || id === props.page.id) return
  if (event.dataTransfer) event.dataTransfer.dropEffect = 'move'

  const el = rowRef.value
  if (!el) return
  const rect = el.getBoundingClientRect()
  const ratio = (event.clientY - rect.top) / rect.height
  dropPos.value = ratio < 0.25 ? 'before' : ratio > 0.75 ? 'after' : 'inside'
}

function onDragLeave(event: DragEvent) {
  // Ignore leave events fired while moving onto our own child elements.
  const el = rowRef.value
  if (el && event.relatedTarget instanceof Node && el.contains(event.relatedTarget)) return
  dropPos.value = null
}

async function onDrop() {
  const pos = dropPos.value
  const id = dnd.draggingId.value
  dropPos.value = null
  if (!pos || id === null || id === props.page.id) return
  await dnd.move(id, props.page.id, pos)
}
</script>