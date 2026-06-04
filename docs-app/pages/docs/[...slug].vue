<script setup lang="ts">
import type {Page} from "~/types/page";

const route = useRoute()
const slug = computed(() => {
  const s = route.params.slug
  return (Array.isArray(s) ? s.join('/') : s) ?? ''
})

// The URL stays slug-based for readability, but the API is keyed by id. Resolve
// the id from the cached pages tree (fetched in app.vue / PageList) so the
// request can use it.
const {data: pages} = useNuxtData<Page[]>('pages')

// Walk the slug path segments down the tree (`parent/child` -> nested page).
function findIdByPath(list: Page[] | null | undefined, segments: string[]): number | undefined {
  const [head, ...rest] = segments
  const match = (list ?? []).find(p => p.slug === head)
  if (!match) return undefined
  return rest.length ? findIdByPath(match.pages, rest) : match.id
}

const pageId = computed(() =>
    slug.value ? findIdByPath(pages.value, slug.value.split('/')) : undefined,
)

const {data: page} = await useAsyncData(
    () => `docs-${slug.value}`,
    () => pageId.value !== undefined
        ? $fetch<Page>(`/api/page/${pageId.value}`)
        : Promise.resolve(null),
    {watch: [slug, pageId]}
)

// Nuxt Content TOC links, fed straight into <UContentToc> (which provides its
// own scroll-spy highlighting).
const links = computed(() => page.value?.toc?.links ?? [])

useSeoMeta({
  title: () => page.value?.title ? `${page.value.title} — DevDocs` : 'DevDocs',
})

useHead(computed(() => ({
  style: (page.value as any)?.highlightCSS
      ? [{id: 'shiki-highlight', children: (page.value as any).highlightCSS}]
      : [],
})))
</script>

<template>

  <UPage :ui="{
    root: 'lg:gap-0 h-full min-h-0 lg:grid-rows-1 mx-auto w-full max-w-[96rem]',
    left: 'lg:col-span-3',
    center: 'lg:col-span-5',
  }">
    <template #left>
      <!-- Sidebar: separated by a hard border (no background). -->
      <PageList
          class="hidden lg:block lg:px-3 lg:border-r lg:border-gray-200 lg:dark:border-gray-800"
      />
    </template>

    <UPageBody class="mt-0 pb-0 h-full min-h-0">
      <!-- Main content: a lighter surface than the rail (owns its own header +
           scroll region). -->
      <div
          class="h-full min-h-0 flex flex-col overflow-hidden bg-white dark:bg-gray-950"
      >
        <!-- Content fills the column (no width cap); Page owns its own header,
             padding and scroll region, so just bound the height here. -->
        <div class="flex-1 min-h-0 w-full">
          <Page :page="page" :slug="slug"/>
        </div>
      </div>
    </UPageBody>

    <template #right>
      <!-- Table of Contents -->
      <PageTocs :links="links"/>
    </template>
  </UPage>
</template>