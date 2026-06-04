<script setup lang="ts">
import type {PageTocLinks} from "~/types/page";

// Nuxt Content TOC links (`page.toc.links`); UContentToc handles the
// scroll-spy highlighting itself.
defineProps<{
  links: PageTocLinks
}>()
</script>

<template>
  <!-- Wrapper takes the (stretching) grid cell; the card sizes to the list and
       is pushed to the side, suspended from a "string" so it reads as hanging
       from the top. -->
  <aside class="hidden xl:block pr-4 pl-3">
    <div v-if="links.length" class="ml-auto w-64 flex flex-col">
      <!-- Two cords the card hangs from, dropped low. -->
      <div class="flex justify-between px-6">
        <span v-for="i in 2" :key="i" class="flex flex-col items-center">
          <span class="w-px h-12 bg-gradient-to-b from-transparent to-gray-300 dark:to-gray-600"/>
          <span class="-mt-1 size-2 rounded-full bg-primary-500 ring-4 ring-primary-500/15"/>
        </span>
      </div>

      <UCard
          class="-mt-0.5 w-full origin-top rounded-2xl shadow-xl shadow-gray-400/20 dark:shadow-black/50 ring-1 ring-gray-200 dark:ring-gray-800 transition-transform duration-300 hover:-translate-y-0.5"
          :ui="{ body: 'p-3 sm:p-3' }"
      >
        <UContentToc
            title="On this page"
            highlight
            highlight-color="primary"
            highlight-variant="circuit"
            :links="links"
            :ui="{
              root: 'border-0 bg-transparent',
              container: 'pt-0 sm:pt-0 lg:pt-0',
              title: 'uppercase font-mono tracking-wider mb-2',
              link: 'text-gray-600 dark:text-gray-400 hover:text-primary-600 dark:hover:text-primary-400 py-0.5',
              linkText: 'whitespace-normal',
            }"
        />
      </UCard>
    </div>
  </aside>
</template>