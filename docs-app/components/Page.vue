<script setup lang="ts">
import type {Page} from "~/types/page";

const props = defineProps<{
  page?: Page,
  slug: string,
}>()

const breadcrumbs = computed(() => [
  {label: 'Home', to: '/', icon: 'i-heroicons-home'},
  {label: 'Docs', to: '/docs/getting-started'},
  {label: props.page?.title ?? props.slug},
])

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
</template>

<style scoped>
/* Keep anchored headings clear of the sticky header when jumped to. */
.prose :deep(:is(h1, h2, h3, h4, h5, h6)[id]) {
  scroll-margin-top: 96px;
}
</style>