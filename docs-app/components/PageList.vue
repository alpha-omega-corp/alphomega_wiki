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


</script>