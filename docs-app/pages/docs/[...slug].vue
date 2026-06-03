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
  <div>
    <UContainer>
      <div class="flex gap-10 py-10">
        <!-- Sidebar -->
        <PageList class="hidden lg:block"/>

        <!-- Main content -->
        <Page :page="page" :slug="slug"/>

        <!-- Table of Contents -->
        <PageTocs :links="links"/>
      </div>
    </UContainer>


  </div>
</template>