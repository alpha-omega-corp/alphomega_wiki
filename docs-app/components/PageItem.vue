<template>
  <li>
    <div class="flex items-center gap-0.5 group">
      <UButton
          v-if="hasChildren"
          :icon="open ? 'i-lucide-chevron-down' : 'i-lucide-chevron-right'"
          size="xs"
          variant="ghost"
          color="neutral"
          class="shrink-0"
          @click="open = !open"
      />
      <span v-else class="size-6 shrink-0"/>

      <NuxtLink
          :to="`/docs/${page.slug}`"
          class="flex flex-1 items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors min-w-0"
          :class="isActive
          ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400 font-medium'
          : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-800'"
      >
        <span class="truncate">{{ page.title }}</span>
      </NuxtLink>

      <UButton
          icon="i-lucide-plus"
          size="xs"
          variant="ghost"
          color="neutral"
          class="shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
          @click="$emit('create', page.id)"
      />
    </div>

    <ul
        v-if="hasChildren"
        v-show="open"
        class="mt-0.5 ml-3 space-y-0.5 border-l border-gray-200 dark:border-gray-700 pl-2"
    >
      <PageItem
          v-for="child in page.pages"
          :key="child.id"
          :page="child"
          @create="$emit('create', $event)"
      />
    </ul>
  </li>
</template>

<script setup lang="ts">
import type {Page} from "~/types/page";

const props = defineProps<{ page: Page }>()
defineEmits<{ create: [id: number] }>()

const route = useRoute()
const open = ref(false)

const hasChildren = computed(() => !!props.page.pages?.length)
const isActive = computed(() => route.path === `/docs/${props.page.slug}`)

function pageInSubtree(page: Page, path: string): boolean {
  if (`/docs/${page.slug}` === path) return true
  return (page.pages ?? []).some(p => pageInSubtree(p, path))
}

const isActiveSubtree = computed(() => pageInSubtree(props.page, route.path))

watchEffect(() => {
  if (isActiveSubtree.value) open.value = true
})
</script>