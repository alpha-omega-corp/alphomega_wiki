<script setup lang="ts">
import Editor from "~/components/editor/Editor.vue"

const route = useRoute()
const slug = computed(() => {
  const s = route.params.slug
  return Array.isArray(s) ? s.join('/') : s
})

const {data: page} = await useAsyncData(
    `docs-${slug.value}`,
    () => $fetch(`/api/page/${slug.value}`),
    {watch: [slug]}
)

const breadcrumbs = computed(() => [
  {label: 'Home', to: '/', icon: 'i-heroicons-home'},
  {label: 'Docs', to: '/docs/getting-started'},
  {label: page.value?.title ?? page.value?.slug},
])

const editorOpen = ref(false)
defineShortcuts({
  o: () => (editorOpen.value = !editorOpen.value)
})

type TocLink = { id: string; text: string; children?: TocLink[] }
type FlatLink = { id: string; text: string; indent: number }

function flattenToc(links: TocLink[], indent = 0): FlatLink[] {
  return links.flatMap(link => [
    {id: link.id, text: link.text, indent},
    ...flattenToc(link.children ?? [], indent + 1),
  ])
}


</script>

<template>
  <main class="flex-1 min-w-0">
    <template v-if="page">
      <!-- Breadcrumb -->
      <div class="flex justify-between items-center mb-6">
        <UBreadcrumb :items="breadcrumbs"/>
        <UDrawer v-model:open="editorOpen" handle-only should-scale-background set-background-color-on-scale>
          <UButton label="Edit" color="primary" variant="subtle" icon="i-lucide-eye"/>

          <template #content>
            <UContainer class="overflow-y-scroll">
              <Editor :value="page.content"></Editor>
            </UContainer>
          </template>
        </UDrawer>
      </div>


      <article class="prose dark:prose-invert max-w-none overflow-y-scroll">
        <ContentRenderer :key="page.slug" :value="page.body" class="mt-4"/>
      </article>

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

</style>