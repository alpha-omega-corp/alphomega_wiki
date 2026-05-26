<template>
  <NuxtLink
    :to="to"
    class="flex items-center gap-2 px-3 py-1.5 rounded-md text-sm transition-colors"
    :class="isActive
      ? 'bg-primary-50 dark:bg-primary-900/20 text-primary-600 dark:text-primary-400 font-medium'
      : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-800'"
  >
    <UIcon name="i-heroicons-document-text" class="size-4 flex-shrink-0" />
    {{ label }}
  </NuxtLink>
</template>

<script setup lang="ts">
const props = defineProps<{ to: string; label: string }>()
const route = useRoute()
const isActive = computed(() => route.path === props.to || route.path.startsWith(props.to + '/'))
</script>
