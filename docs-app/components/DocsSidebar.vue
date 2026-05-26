<template>
  <nav class="w-64 flex-shrink-0">
    <div v-if="status === 'pending'" class="space-y-4 px-3">
      <div v-for="i in 3" :key="i" class="space-y-1.5">
        <USkeleton class="h-3 w-2/3 rounded" />
        <USkeleton v-for="j in 3" :key="j" class="h-7 rounded-md" />
      </div>
    </div>

    <div v-else-if="status === 'error'" class="px-3 py-2 text-sm text-red-500">
      Failed to load navigation
    </div>

    <ul v-else class="space-y-0.5">
      <SidebarItem
        v-for="page in pages"
        :key="page.id"
        :page="page"
        @create="openModal"
      />
    </ul>

    <CreateArticleModal
      v-model:open="modalOpen"
      v-model:parent-id="selectedParentId"
      :pages="allPages"
      @created="refresh()"
    />
  </nav>
</template>

<script setup lang="ts">
import type { NavPage } from '~/types/navigation'

const { data: pages, status, refresh } = await useAsyncData(
  'navigation',
  () => $fetch<NavPage[]>('/api/navigation'),
  { default: () => [] as NavPage[] },
)

function flattenPages(pages: NavPage[]): NavPage[] {
  return pages.flatMap(p => [p, ...flattenPages(p.pages ?? [])])
}

const allPages = computed<NavPage[]>(() => flattenPages(pages.value))

const modalOpen = ref(false)
const selectedParentId = ref<number | null>(null)

function openModal(id: number) {
  selectedParentId.value = id
  modalOpen.value = true
}
</script>