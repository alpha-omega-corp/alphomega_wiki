<template>
  <div>
    <!-- Hero -->
    <section class="py-24 text-center">
      <UContainer>
        <UBadge label="v2.0 Released" variant="subtle" class="mb-6" />
        <h1 class="text-5xl font-bold tracking-tight text-gray-900 dark:text-white mb-6">
          Build with
          <span class="text-primary-500">confidence</span>
        </h1>
        <p class="text-xl text-gray-500 dark:text-gray-400 max-w-2xl mx-auto mb-10">
          A complete developer platform with comprehensive documentation, live examples,
          and powerful editing tools to accelerate your workflow.
        </p>
        <div class="flex items-center justify-center gap-4">
          <UButton label="Get Started" to="/docs/getting-started" size="xl" />
          <UButton label="View on GitHub" icon="i-heroicons-arrow-top-right-on-square" variant="ghost" color="neutral" size="xl" />
        </div>
      </UContainer>
    </section>

    <!-- Feature Cards -->
    <section class="py-16 bg-white dark:bg-gray-900">
      <UContainer>
        <div class="text-center mb-12">
          <h2 class="text-3xl font-bold text-gray-900 dark:text-white mb-4">Everything you need</h2>
          <p class="text-gray-500 dark:text-gray-400">Built for modern development teams</p>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <UCard v-for="feature in features" :key="feature.title">
            <div class="flex flex-col gap-3">
              <div class="size-10 rounded-lg bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center">
                <UIcon :name="feature.icon" class="size-5 text-primary-500" />
              </div>
              <h3 class="font-semibold text-gray-900 dark:text-white">{{ feature.title }}</h3>
              <p class="text-sm text-gray-500 dark:text-gray-400">{{ feature.description }}</p>
            </div>
          </UCard>
        </div>
      </UContainer>
    </section>

    <!-- Quick Start -->
    <section class="py-16">
      <UContainer>
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          <div>
            <h2 class="text-3xl font-bold text-gray-900 dark:text-white mb-4">Quick start in seconds</h2>
            <p class="text-gray-500 dark:text-gray-400 mb-6">
              Get up and running with a single command. Our CLI handles everything
              from project scaffolding to dependency management.
            </p>
            <UButton label="Read the docs" to="/docs/getting-started" />
          </div>
          <div class="bg-gray-900 rounded-xl p-6 font-mono text-sm">
            <div class="flex items-center gap-2 mb-4">
              <div class="size-3 rounded-full bg-red-500" />
              <div class="size-3 rounded-full bg-yellow-500" />
              <div class="size-3 rounded-full bg-green-500" />
            </div>
            <p class="text-gray-400">$ <span class="text-green-400">npx</span> <span class="text-white">create-my-app my-project</span></p>
            <p class="text-gray-500 mt-1">$ <span class="text-green-400">cd</span> <span class="text-white">my-project</span></p>
            <p class="text-gray-500 mt-1">$ <span class="text-green-400">npm</span> <span class="text-white">run dev</span></p>
            <p class="text-green-400 mt-3">✓ Server running at http://localhost:3000</p>
          </div>
        </div>
      </UContainer>
    </section>

    <!-- Doc Categories -->
    <section class="py-16 bg-white dark:bg-gray-900">
      <UContainer>
        <h2 class="text-3xl font-bold text-gray-900 dark:text-white mb-8 text-center">Browse documentation</h2>
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
const features = [
  {
    icon: 'i-heroicons-document-text',
    title: 'Rich Documentation',
    description: 'Markdown-powered docs with syntax highlighting, code examples, and live previews.',
  },
  {
    icon: 'i-heroicons-pencil-square',
    title: 'Built-in Editor',
    description: 'Write and edit content directly in the browser with our WYSIWYG editor.',
  },
  {
    icon: 'i-heroicons-squares-2x2',
    title: 'Component Library',
    description: 'Over 50 accessible, customizable UI components ready to drop into your project.',
  },
]

const docCategories = [
  { label: 'Getting Started', icon: 'i-heroicons-rocket-launch', to: '/docs/getting-started', description: 'Installation, setup, and your first project' },
  { label: 'Configuration', icon: 'i-heroicons-cog-6-tooth', to: '/docs/configuration', description: 'Environment variables, config options' },
  { label: 'API Reference', icon: 'i-heroicons-code-bracket', to: '/docs/api', description: 'Full API docs and type definitions' },
  { label: 'Deployment', icon: 'i-heroicons-cloud-arrow-up', to: '/docs/deployment', description: 'Deploy to any platform in minutes' },
]
</script>
