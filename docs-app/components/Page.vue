<script setup lang="ts">
import type {BreadcrumbItem} from "@nuxt/ui";
import type {Page} from "~/types/page";
import WCard from "~/components/content/WCard.vue";

// Render `::card` blocks with our component instead of Nuxt UI's default
// ProseCard (passed components override the built-in prose map).
const mdcComponents = {card: WCard}

const props = defineProps<{
  page?: Page | null,
  slug: string,
}>()

// Editing is gated behind login (the server also enforces this on every
// mutation). Logged-out users get a read-only view.
const {loggedIn, login} = useAuth()
const route = useRoute()

// Shared page tree (fetched in PageList / app.vue) used to resolve each slug
// segment to its page title + link.
const {data: pageTree} = useNuxtData<Page[]>('pages')

// Build a crumb per slug segment so nested pages show their full path, e.g.
// Home › Authentication › SSO / OAuth — each ancestor linking to its own page.
const breadcrumbs = computed<BreadcrumbItem[]>(() => {
  const crumbs: BreadcrumbItem[] = [{label: 'Home', to: '/', icon: 'i-heroicons-home'}]

  const segments = props.slug.split('/').filter(Boolean)
  let level = pageTree.value ?? []
  let path = ''

  segments.forEach((segment, index) => {
    path = path ? `${path}/${segment}` : segment
    const match = level.find(p => p.slug === segment)
    const isLast = index === segments.length - 1

    crumbs.push({
      // Prefer the page's real title; fall back to the live page prop for the
      // current (last) crumb, then to the raw slug segment.
      label: match?.title ?? (isLast ? props.page?.title : undefined) ?? segment,
      // The current page is plain text (no link).
      to: isLast ? undefined : `/docs/${path}`,
    })

    level = match?.pages ?? []
  })

  return crumbs
})

// A page with no markdown yet (e.g. just created) is opened straight into the
// editor instead of the empty read-only view + drawer.
const hasContent = computed(() => !!props.page?.content?.trim())

const editorOpen = ref(false)
defineShortcuts({
  o: () => (editorOpen.value = !editorOpen.value)
})

// After saving in the editor, re-fetch the page so the rendered article reflects
// the new content, and close the drawer.
async function onSaved() {
  await refreshNuxtData(`docs-${props.slug}`)
  editorOpen.value = false
}
</script>

<template>
  <main class="flex flex-col h-full min-h-0 min-w-0 pt-6">
    <template v-if="page">
      <!-- Breadcrumb: fixed header band, kept out of the scroll region. Spans
           the full content width with its own padding. -->
      <div
          class="flex justify-between items-center gap-4 flex-shrink-0 -mt-6 p-5  border-gray-200 dark:border-gray-800">
        <UBreadcrumb :items="breadcrumbs"/>
        <UDrawer v-if="hasContent && loggedIn" v-model:open="editorOpen" handle-only should-scale-background
                 set-background-color-on-scale>
          <UButton label="Edit" color="secondary" variant="subtle" icon="i-lucide-eye"/>

          <template #content>
            <UContainer class="overflow-y-scroll">
              <Editor :value="page.content" :page-id="page.id" @saved="onSaved"></Editor>
            </UContainer>
          </template>
        </UDrawer>
      </div>

      <!-- Scrolling content region (everything below the breadcrumb). Padding
           lives here so the scrollbar sits in the gutter, clear of the text; the
           right margin keeps the scrollbar off the table of contents. -->
      <div class="flex-1 min-h-0 overflow-y-auto px-4 sm:px-5 lg:px-6 mr-3">
        <!-- Empty page: edit inline so the user can start writing right away. -->
        <Editor
            v-if="!hasContent && loggedIn"
            :key="page.id"
            :value="page.content"
            :page-id="page.id"
            @saved="onSaved"
        />

        <!-- Empty page, logged out: read-only prompt to sign in. -->
        <div v-else-if="!hasContent" class="py-16 text-center">
          <UIcon name="i-lucide-file-text" class="size-10 text-gray-400 mx-auto mb-3"/>
          <p class="text-gray-500 mb-4">This page has no content yet.</p>
          <UButton size="sm" icon="i-lucide-log-in" label="Sign in to edit" @click="login(route.fullPath)"/>
        </div>

        <template v-else>
          <article class="prose dark:prose-invert max-w-none mt-6">
            <ContentRenderer :key="slug" :value="page" :components="mdcComponents"/>
          </article>
        </template>
      </div>
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
</template>

<style scoped>
/* Keep anchored headings clear of the sticky header when jumped to. */
.prose :deep(:is(h1, h2, h3, h4, h5, h6)[id]) {
  scroll-margin-top: 96px;
}
</style>