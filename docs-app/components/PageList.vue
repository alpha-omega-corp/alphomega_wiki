<template>
  <nav class="flex-shrink-0 p-5 h-full overflow-y-auto">

    <!-- The current user's favorited pages, resolved to their slug paths. -->
    <div v-if="favoritePages.length" class="mb-4">
      <p class="px-2 mb-1 text-xs font-semibold uppercase tracking-wider text-gray-400 dark:text-gray-500">
        Favorites
      </p>
      <ul class="space-y-0.5">
        <li
            v-for="fav in favoritePages"
            :key="fav.page.id"
            class="group flex items-center gap-0.5 rounded-md"
        >
          <NuxtLink
              :to="`/docs/${fav.slug}`"
              class="flex flex-1 items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors min-w-0"
              :class="route.path === `/docs/${fav.slug}`
                ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400 font-medium'
                : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-800'"
          >
            <span class="flex size-4 shrink-0 items-center justify-center">
              <UIcon
                  v-if="fav.page.icon && isIcon(fav.page.icon)"
                  :name="fav.page.icon"
                  class="size-4 text-gray-500 dark:text-gray-400"
              />
              <span v-else-if="fav.page.icon" class="text-sm leading-none">{{ fav.page.icon }}</span>
              <UIcon v-else name="i-lucide-star" class="size-4 text-amber-400"/>
            </span>
            <span class="truncate">{{ fav.page.title }}</span>
          </NuxtLink>

          <UButton
              icon="i-lucide-star-off"
              size="xs"
              variant="ghost"
              color="neutral"
              class="shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
              aria-label="Remove from favorites"
              @click="toggleFavorite(fav.page)"
          />
        </li>
      </ul>
    </div>

    <UContentSearchButton :collapsed="false" class="mb-4 w-full"/>

    <ul class="space-y-0.5">
      <PageItem
          v-for="page in pages"
          :key="page.id"
          :page="page"
          @create="openModal"
      />
    </ul>

    <div v-if="status === 'success' && !pages?.length" class="mt-4 flex flex-col items-center gap-2 rounded-lg border border-dashed border-gray-200 dark:border-gray-700 px-4 py-6 text-center">
      <UIcon name="i-lucide-file-plus" class="size-7 text-gray-300 dark:text-gray-600" />
      <p class="text-xs text-gray-400 dark:text-gray-500">No pages yet</p>
      <UButton
          size="xs"
          icon="i-lucide-plus"
          label="Create page"
          @click="selectedParentId = null; createPageModal = true"
      />
    </div>

    <CreatePage
        v-model:open="createPageModal"
        v-model:parent-id="selectedParentId"
    />
  </nav>
</template>

<script setup lang="ts">
import type {Page} from "~/types/page";
import {type DropPosition, type PageDnd, PageDndKey} from "~/composables/pageDnd";

const {data: pages, status} = await useAsyncData(
    'pages',
    () => $fetch<Page[]>('/api/pages')
)

// Load the current user's favorites once here (PageList owns the tree); every
// PageItem then reads the shared 'favorites' cache via useFavorites().
loadFavorites()

const route = useRoute()
const {favorites, toggle: toggleFavorite} = useFavorites()

function isIcon(value: string) {
  return /^i-[a-z0-9]+[:-]/.test(value) || value.includes(':')
}

// Locate a page in the tree and build its full slug path (including ancestors).
function findWithSlug(siblings: Page[], id: number, base = ''): { page: Page; slug: string } | null {
  for (const node of siblings) {
    const slug = base ? `${base}/${node.slug}` : node.slug
    if (node.id === id) return {page: node, slug}
    if (node.pages?.length) {
      const found = findWithSlug(node.pages, id, slug)
      if (found) return found
    }
  }
  return null
}

// Favorited pages resolved to tree nodes + slug paths, in backend order (newest
// first). Skips any favorite whose page no longer exists in the tree.
const favoritePages = computed(() => {
  const roots = pages.value ?? []
  return favorites.value
      .map(f => findWithSlug(roots, f.page_id))
      .filter((x): x is { page: Page; slug: string } => x !== null)
})

const createPageModal = ref(false)
const selectedParentId = ref<number | null>(null)

function openModal(id: number) {
  selectedParentId.value = id
  createPageModal.value = true
}

// --- Drag & drop -----------------------------------------------------------
// PageList owns the page tree, so it also owns the move logic; PageItem nodes
// only report "page X was dropped <before|after|inside> page Y".

const draggingId = ref<number | null>(null)

type Ctx = { node: Page; siblings: Page[]; index: number; parentId: number | null }

// Locate a page in the tree along with the array it lives in and its parent id.
function findCtx(siblings: Page[], id: number, parentId: number | null = null): Ctx | null {
  for (let i = 0; i < siblings.length; i++) {
    const node = siblings[i]
    if (!node) continue
    if (node.id === id) return {node, siblings, index: i, parentId}
    if (node.pages?.length) {
      const found = findCtx(node.pages, id, node.id)
      if (found) return found
    }
  }
  return null
}

function inSubtree(node: Page, id: number): boolean {
  if (node.id === id) return true
  return (node.pages ?? []).some(child => inSubtree(child, id))
}

async function move(draggedId: number, targetId: number, position: DropPosition) {
  const roots = pages.value
  if (!roots || draggedId === targetId) return

  const dragged = findCtx(roots, draggedId)
  const target = findCtx(roots, targetId)
  if (!dragged || !target) return

  // Resolve the destination sibling list + its parent id.
  let destParentId: number | null
  let destSiblings: Page[]
  if (position === 'inside') {
    destParentId = targetId
    if (!target.node.pages) target.node.pages = []
    destSiblings = target.node.pages
  } else {
    destParentId = target.parentId
    destSiblings = target.siblings
  }

  // Never move a page into itself or one of its own descendants.
  if (destParentId !== null && inSubtree(dragged.node, destParentId)) return

  // Detach from the current location, then insert at the destination. (The
  // target's index may shift when moving within the same sibling list, so
  // recompute it after the removal.)
  dragged.siblings.splice(dragged.index, 1)
  let insertAt: number
  if (position === 'inside') {
    insertAt = destSiblings.length
  } else {
    const ti = destSiblings.indexOf(target.node)
    insertAt = position === 'before' ? ti : ti + 1
  }
  destSiblings.splice(insertAt, 0, dragged.node)

  // Renumber the affected sibling lists; persist every page whose parent or
  // position changed. (Renumbering both lists is harmless when they're equal.)
  const changes = new Map<number, { parent_id: number | null; position: number }>()
  const renumber = (siblings: Page[], parentId: number | null) =>
      siblings.forEach((p, i) => changes.set(p.id, {parent_id: parentId, position: i}))
  renumber(dragged.siblings, dragged.parentId)
  renumber(destSiblings, destParentId)

  // Reveal the new children when nesting into a collapsed parent.
  if (position === 'inside') expanded.value[targetId] = true

  try {
    await Promise.all([...changes].map(([id, body]) =>
        $fetch(`/api/page/${id}`, {method: 'PUT', body}),
    ))
  } finally {
    // Re-sync with the server's authoritative ordering (and roll back on error).
    await refreshNuxtData('pages')
  }
}

const dnd: PageDnd = {
  draggingId,
  start: (id) => (draggingId.value = id),
  end: () => (draggingId.value = null),
  move,
}
provide(PageDndKey, dnd)

// Persist the expanded-tree state (shared via useState in PageItem) to
// localStorage so it survives full page reloads, not just navigations.
const STORAGE_KEY = 'pagelist-expanded'
const expanded = useState<Record<number, boolean>>('pagelist-expanded', () => ({}))

onMounted(() => {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) {
      // Restore saved state, but let anything opened during this load (e.g. the
      // active page's subtree) take precedence.
      expanded.value = {...JSON.parse(saved), ...expanded.value}
    }
  } catch {
    // Ignore unavailable/corrupt storage.
  }

  watch(expanded, (value) => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(value))
    } catch {
      // Ignore quota/availability errors.
    }
  }, {deep: true})
})


</script>