<script setup lang="ts">
import type {PageTocLinks} from "~/types/page";

// Nuxt Content TOC links (`page.toc.links`); UContentToc handles the
// scroll-spy highlighting itself.
defineProps<{
  links: PageTocLinks
}>()
</script>

<template>
  <aside class="hidden xl:block w-56 flex-shrink-0">
    <UContentToc
        v-if="links.length"
        highlight
        highlight-color="primary"
        highlight-variant="circuit"
        :links="links"
        class="sticky top-24"
    />
  </aside>
</template>