<template>
  <div>
    <!-- Hero -->
    <section class="px-6 pt-20 pb-16 sm:pt-28 sm:pb-20">
      <div class="mx-auto max-w-2xl text-center">
        <p class="text-sm font-medium uppercase tracking-widest text-primary-500">
          DevDocs
        </p>
        <h1 class="mt-4 text-4xl sm:text-5xl font-bold tracking-tight text-gray-900 dark:text-white">
          Everything you need, <br class="hidden sm:block" />beautifully documented
        </h1>
        <p class="mt-6 text-lg leading-relaxed text-gray-500 dark:text-gray-400">
          Guides, references, and deployment recipes — written to get you from
          zero to shipped without the noise.
        </p>
        <div class="mt-10 flex items-center justify-center gap-3">
          <UButton
              to="/docs/getting-started"
              size="lg"
              trailing-icon="i-heroicons-arrow-right"
          >
            Get started
          </UButton>
          <UButton
              to="/docs/api"
              size="lg"
              color="neutral"
              variant="subtle"
          >
            API reference
          </UButton>
        </div>
      </div>
    </section>

    <!-- Doc Categories -->
    <section class="px-6 pb-24">
      <UContainer>
        <h2 class="mb-8 text-sm font-semibold uppercase tracking-wider text-gray-400 dark:text-gray-500">
          Explore the docs
        </h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          <NuxtLink
              v-for="cat in docCategories"
              :key="cat.to"
              :to="cat.to"
              class="group"
          >
            <UCard
                class="h-full transition-all duration-200 hover:-translate-y-1 hover:ring-2 hover:ring-primary-500"
                :ui="{ body: 'p-6 sm:p-7' }"
            >
              <div class="flex flex-col gap-4">
                <div class="flex size-11 items-center justify-center rounded-xl bg-primary-50 dark:bg-primary-950 text-primary-500">
                  <UIcon :name="cat.icon" class="size-5" />
                </div>
                <div class="flex flex-col gap-1.5">
                  <h3 class="font-semibold text-gray-900 dark:text-white group-hover:text-primary-500 transition-colors">
                    {{ cat.label }}
                  </h3>
                  <p class="text-sm leading-relaxed text-gray-500 dark:text-gray-400">
                    {{ cat.description }}
                  </p>
                </div>
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
