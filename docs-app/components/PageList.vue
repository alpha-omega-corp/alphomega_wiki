<template>
  <nav class="flex-shrink-0 p-5 h-full overflow-y-auto">

    <UContentSearchButton :collapsed="false" class="mb-4 w-full"/>

    <ul class="space-y-0.5">
      <PageItem
          v-for="page in pages"
          :key="page.id"
          :page="page"
          @create="openModal"
      />
    </ul>

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