<template>
  <UApp>
    <NuxtLayout>
      <NuxtPage/>
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
import type {ContentNavigationItem} from '@nuxt/content'
import type {Page} from '~/types/page'


// Shares the same cache key as DocsSidebar so the tree is only fetched once.
const {data: pages} = await useAsyncData(
    'pages',
    () => $fetch<Page[]>('/api/pages')
)

// Flatten the page tree: every page is a real doc, so each becomes a
// searchable leaf (titles/content are matched by the command palette).
type SearchFile = { id: string; title: string; titles: string[]; level: number; content: string }

// `basePath`/`titles` accumulate the ancestor chain so nested pages get their
// full slug path (e.g. `parent/child`) — the tree nodes only carry their own
// slug segment, so a bare `/docs/${page.slug}` would point at the wrong URL.
function flatten(items: Page[], basePath = '', titles: string[] = [], acc: SearchFile[] = []): SearchFile[] {
  for (const page of items) {
    const fullSlug = basePath ? `${basePath}/${page.slug}` : page.slug
    acc.push({
      id: `/docs/${fullSlug}`,
      title: page.title,
      titles,
      level: 1,
      content: page.content ?? '',
    })
    if (page.pages?.length) flatten(page.pages, fullSlug, [...titles, page.title], acc)
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