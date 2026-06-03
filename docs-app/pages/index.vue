<template>
  <div>
    <!-- Doc Categories -->
    <section class="py-16 bg-white dark:bg-gray-900">
      <UContainer>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <NuxtLink
              v-for="cat in docCategories"
              :key="cat.to"
              :to="cat.to"
              class="group"
          >
            <UCard class="h-full hover:ring-2 hover:ring-primary-500 transition-all">
              <div class="flex flex-col gap-2">
                <UIcon :name="cat.icon" class="size-6 text-primary-500" />
                <h3 class="font-semibold text-gray-900 dark:text-white group-hover:text-primary-500 transition-colors">{{ cat.label }}</h3>
                <p class="text-xs text-gray-500">{{ cat.description }}</p>
              </div>
            </UCard>
          </NuxtLink>
        </div>
      </UContainer>
    </section>
  </div>
</template>

<script setup lang="ts">
const docCategories = [
  { label: 'Getting Started', icon: 'i-heroicons-rocket-launch', to: '/docs/getting-started', description: 'Installation, setup, and your first project' },
  { label: 'Configuration', icon: 'i-heroicons-cog-6-tooth', to: '/docs/configuration', description: 'Environment variables, config options' },
  { label: 'API Reference', icon: 'i-heroicons-code-bracket', to: '/docs/api', description: 'Full API docs and type definitions' },
  { label: 'Deployment', icon: 'i-heroicons-cloud-arrow-up', to: '/docs/deployment', description: 'Deploy to any platform in minutes' },
]
</script>
