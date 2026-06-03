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
    {id: link.id, text: link.text, indent},
    ...flattenToc(link.children ?? [], indent + 1),
  ])
}

const {data: page, status} = await useAsyncData(
    () => `docs-${slug.value}`,
    () => $fetch(`/api/page/${slug.value}`),
    {watch: [slug]}
)


const tocLinks = computed<FlatLink[]>(() =>
    flattenToc(page.value?.toc?.links ?? [])
)

// Scroll-spy: highlight every TOC entry whose heading is currently visible in
// the viewport, so the reader can see how much of the page is on screen and how
// deep into the document they are.
const activeIds = ref<string[]>([])
let observer: IntersectionObserver | null = null

function setupScrollSpy() {
  observer?.disconnect()
  if (!import.meta.client) return

  const ids = tocLinks.value.map(l => l.id)
  const headings = ids
      .map(id => document.getElementById(id))
      .filter((el): el is HTMLElement => el !== null)

  if (!headings.length) {
    activeIds.value = []
    return
  }

  const visible = new Set<string>()

  observer = new IntersectionObserver((entries) => {
    for (const entry of entries) {
      if (entry.isIntersecting) visible.add(entry.target.id)
      else visible.delete(entry.target.id)
    }

    // Highlight all headings currently on screen (in document order).
    const onScreen = ids.filter(id => visible.has(id))
    if (onScreen.length) {
      activeIds.value = onScreen
      return
    }

    // Between headings (a long section fills the viewport): keep the last
    // heading we scrolled past highlighted so something stays active.
    const passed = headings.filter(el => el.getBoundingClientRect().top < 96)
    activeIds.value = passed.length ? [passed[passed.length - 1].id] : [ids[0]]
  }, {
    // Count a heading as visible once it clears the sticky header (~96px) and
    // until it leaves the bottom of the viewport.
    rootMargin: '-96px 0px 0px 0px',
    threshold: 0,
  })

  headings.forEach(el => observer!.observe(el))
}

onMounted(setupScrollSpy)
onBeforeUnmount(() => observer?.disconnect())
// Re-wire the observer when navigating to a different page (this component is
// reused across /docs/* routes, so the headings change without a remount).
watch(tocLinks, () => nextTick(setupScrollSpy))

const breadcrumbs = computed(() => [
  {label: 'Home', to: '/', icon: 'i-heroicons-home'},
  {label: 'Docs', to: '/docs/getting-started'},
  {label: page.value?.title ?? slug.value},
])

const editorOpen = ref(false)
defineShortcuts({
  o: () => (editorOpen.value = !editorOpen.value)
})

// After saving in the editor, re-fetch the page so the rendered article reflects
// the new content, and close the drawer.
async function onSaved() {
  await refreshNuxtData(`docs-${slug.value}`)
  editorOpen.value = false
}

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
        <main class="flex-1 min-w-0">
          <template v-if="page">
            <!-- Breadcrumb -->
            <div class="flex justify-between items-center mb-6">
              <UBreadcrumb :items="breadcrumbs"/>
              <UDrawer v-model:open="editorOpen" handle-only should-scale-background set-background-color-on-scale>
                <UButton label="Edit" color="secondary" variant="subtle" icon="i-lucide-eye"/>

                <template #content>
                  <UContainer class="overflow-y-scroll">
                    <Editor :value="page.content" :slug="slug" @saved="onSaved"></Editor>
                  </UContainer>
                </template>
              </UDrawer>
            </div>


            <article class="prose dark:prose-invert max-w-none overflow-y-scroll">
              <ContentRenderer :key="slug" :value="page" class="mt-4"/>
            </article>

            <!-- Extra scroll space below the content. -->
            <div class="h-[200px]" aria-hidden="true"/>

          </template>

          <template v-else>
            <UCard class="text-center py-12">
              <UIcon name="i-heroicons-document-magnifying-glass" class="size-12 text-gray-400 mx-auto mb-4"/>
              <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">Page not found</h2>
              <p class="text-gray-500 mb-6">The documentation page you're looking for doesn't exist.</p>
              <UButton label="Back to docs" to="/docs/getting-started"/>
            </UCard>
          </template>
        </main>

        <!-- Table of Contents -->
        <aside class="hidden xl:block w-56 flex-shrink-0">

          <div class="w-full sticky top-24">
            <nav v-if="tocLinks.length" class="space-y-1 border-l border-gray-200 dark:border-gray-800">
              <a
                  v-for="link in tocLinks"
                  :key="link.id"
                  :href="`#${link.id}`"
                  :style="`padding-left: ${12 + link.indent * 12}px`"
                  class="-ml-px block border-l-2 py-0.5 text-sm transition-colors truncate"
                  :class="activeIds.includes(link.id)
                    ? 'border-primary-500 text-primary-600 dark:text-primary-400 font-medium'
                    : 'border-transparent text-gray-500 hover:text-gray-900 dark:hover:text-white'"
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

<style scoped>
/* Keep anchored headings clear of the sticky header when jumped to. */
.prose :deep(:is(h1, h2, h3, h4, h5, h6)[id]) {
  scroll-margin-top: 96px;
}
</style>