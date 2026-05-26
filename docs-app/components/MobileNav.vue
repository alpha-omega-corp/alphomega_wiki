<template>
  <div>
    <UButton
      icon="i-lucide-menu"
      variant="ghost"
      color="neutral"
      size="sm"
      aria-label="Open navigation"
      @click="open = true"
    />

    <UDrawer v-model:open="open" direction="left" :ui="{ content: 'w-72' }">
      <template #content>
        <div class="flex flex-col h-full">
          <div class="flex items-center justify-between px-4 py-3 border-b border-gray-200 dark:border-gray-800 shrink-0">
            <span class="font-semibold text-gray-900 dark:text-white">Navigation</span>
            <UButton
              icon="i-lucide-x"
              variant="ghost"
              color="neutral"
              size="sm"
              @click="open = false"
            />
          </div>
          <div class="overflow-y-auto flex-1 p-3">
            <DocsSidebar />
          </div>
        </div>
      </template>
    </UDrawer>
  </div>
</template>

<script setup lang="ts">
const open = ref(false)

const route = useRoute()
watch(() => route.path, () => { open.value = false })
</script>
