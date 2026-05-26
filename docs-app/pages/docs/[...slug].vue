
<script setup lang="ts">
const route = useRoute()
const slug = computed(() => {
  const s = route.params.slug
  return Array.isArray(s) ? s.join('/') : s
})

type TocLink = { id: string; text: string; children?: TocLink[] }
type FlatLink = { id: string; text: string; indent: number }

function flattenToc(links: TocLink[], indent = 0): FlatLink[] {
  return links.flatMap(link => [
    { id: link.id, text: link.text, indent },
    ...flattenToc(link.children ?? [], indent + 1),
  ])
}

const { data: page, status } = await useAsyncData(
    () => `docs-${slug.value}`,
    () => $fetch(`/api/page/${slug.value}`),
    { watch: [slug] }
)


const tocLinks = computed<FlatLink[]>(() =>
    flattenToc(page.value?.toc?.links ?? [])
)

const breadcrumbs = computed(() => [
  { label: 'Home', to: '/', icon: 'i-heroicons-home' },
  { label: 'Docs', to: '/docs/getting-started' },
  { label: page.value?.title ?? slug.value },
])

useSeoMeta({
  title: () => page.value?.title ? `${page.value.title} — DevDocs` : 'DevDocs',
})

useHead(computed(() => ({
  style: (page.value as any)?.highlightCSS
    ? [{ id: 'shiki-highlight', children: (page.value as any).highlightCSS }]
    : [],
})))
</script>

<template>
  <div>
    <UContainer>
      <div class="flex gap-10 py-10">
        <!-- Sidebar -->
        <DocsSidebar class="hidden lg:block" />

        <!-- Main content -->
        <main class="flex-1 min-w-0">
          <template v-if="page">
            <!-- Breadcrumb -->
            <div class="flex justify-between items-center mb-6">
              <UBreadcrumb :items="breadcrumbs" />
              <UDrawer handle-only should-scale-background set-background-color-on-scale>
                <UButton label="Edit" color="primary" variant="subtle" icon="i-lucide-eye" />

                <template #content>
                  <UContainer class="overflow-y-scroll">
                    <Editor :value="page.content"></Editor>
                  </UContainer>
                </template>
              </UDrawer>
            </div>


            <article class="prose dark:prose-invert max-w-none overflow-y-scroll">
              <ContentRenderer :key="slug" :value="page" class="mt-4" />
            </article>

          </template>

          <template v-else>
            <UCard class="text-center py-12">
              <UIcon name="i-heroicons-document-magnifying-glass" class="size-12 text-gray-400 mx-auto mb-4" />
              <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">Page not found</h2>
              <p class="text-gray-500 mb-6">The documentation page you're looking for doesn't exist.</p>
              <UButton label="Back to docs" to="/docs/getting-started" />
            </UCard>
          </template>
        </main>

        <!-- Table of Contents -->
        <aside class="hidden xl:block w-56 flex-shrink-0">
          <div class="sticky top-24">
            <p class="text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400 mb-3">On this page</p>
            <template v-if="status === 'pending'">
              <div class="space-y-2">
                <USkeleton v-for="i in 5" :key="i" class="h-4 rounded" :style="`width: ${60 + (i % 3) * 15}%`" />
              </div>
            </template>
            <nav v-else-if="tocLinks.length" class="space-y-1">
              <a
                v-for="link in tocLinks"
                :key="link.id"
                :href="`#${link.id}`"
                :style="link.indent > 0 ? `padding-left: ${link.indent * 12}px` : undefined"
                class="block text-sm text-gray-500 hover:text-gray-900 dark:hover:text-white py-0.5 transition-colors truncate"
              >
                {{ link.text }}
              </a>
            </nav>
          </div>
        </aside>
      </div>
    </UContainer>


  </div>
</template>

