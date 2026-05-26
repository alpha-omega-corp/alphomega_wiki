<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950">
    <!-- Top Navigation Bar -->
    <header class="sticky top-0 z-50 border-b border-gray-200 dark:border-gray-800 bg-white/80 dark:bg-gray-900/80 backdrop-blur">
      <UContainer>
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center gap-3">
            <MobileNav class="lg:hidden" />

            <nav class="hidden md:flex items-center gap-1">
              <UButton
                v-for="item in navItems"
                :key="item.to"
                :to="item.to"
                variant="ghost"
                color="neutral"
                size="lg"
                :label="item.label"
              />
            </nav>
          </div>

          <div class="flex items-center gap-3">
            <UInput
              placeholder="Search docs..."
              icon="i-heroicons-magnifying-glass"
              size="sm"
              class="hidden sm:block w-56"
            />
            <UButton
              icon="i-heroicons-moon"
              variant="ghost"
              color="neutral"
              size="sm"
              @click="colorMode.preference = colorMode.value === 'dark' ? 'light' : 'dark'"
            />
          </div>
        </div>
      </UContainer>
    </header>

    <slot />
  </div>
</template>

<script setup lang="ts">
const colorMode = useColorMode()

const navItems = [
  { label: 'Home', to: '/' },
  { label: 'Documentation', to: '/docs/getting-started' },
]
</script>
