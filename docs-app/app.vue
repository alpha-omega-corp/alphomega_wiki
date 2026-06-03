<template>
  <UApp>
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>

    <ClientOnly>
      <LazyUContentSearch
        :navigation="searchNavigation"
        :files="searchFiles"
        :fuse="{ resultLimit: 20 }"
        :color-mode="false"
      />
    </ClientOnly>
  </UApp>
</template>

<script setup lang="ts">
import type { ContentNavigationItem } from '@nuxt/content'
import type { Page } from '~/types/page'

// Shares the same cache key as DocsSidebar so the tree is only fetched once.
const { data: pages } = await useAsyncData(
  'pages',
  () => $fetch<Page[]>('/api/pages')
)

// Flatten the page tree: every page is a real doc, so each becomes a
// searchable leaf (titles/content are matched by the command palette).
type SearchFile = { id: string; title: string; titles: string[]; level: number; content: string }

function flatten(items: Page[], acc: SearchFile[] = []): SearchFile[] {
  for (const page of items) {
    acc.push({
      id: `/docs/${page.slug}`,
      title: page.title,
      titles: [],
      level: 1,
      content: page.content ?? '',
    })
    if (page.pages?.length) flatten(page.pages, acc)
  }
  return acc
}

const searchFiles = computed<SearchFile[]>(() => flatten(pages.value ?? []))

// A flat navigation list keeps every page (including parents) searchable;
// the `files` above are matched against these paths to populate the palette.
const searchNavigation = computed<ContentNavigationItem[]>(() =>
  searchFiles.value.map(file => ({
    title: file.title,
    path: file.id,
  })),
)
</script>