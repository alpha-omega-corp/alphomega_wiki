<template>
  <nav class="w-64 flex-shrink-0">

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