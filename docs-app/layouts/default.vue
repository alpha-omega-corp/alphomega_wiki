<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-950">
    <!-- Top Navigation Bar -->
    <header
        class="sticky top-0 z-50 border-b border-gray-200 dark:border-gray-800 bg-white/80 dark:bg-gray-900/80 backdrop-blur">
      <UContainer>
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center gap-3">

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

    <!-- App Content -->
    <slot/>

  
    <!-- Footer -->
    <UFooter class="border-t border-gray-200 dark:border-gray-800">
      <template #left>
        <p class="text-sm text-gray-500">
          © {{ new Date().getFullYear() }} DevDocs. All rights reserved.
        </p>
      </template>

      <UButton
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          variant="ghost"
          color="neutral"
          size="sm"
          :label="item.label"
      />

      <template #right>
        <UButton
            icon="i-simple-icons-github"
            to="https://github.com"
            target="_blank"
            variant="ghost"
            color="neutral"
            size="sm"
            aria-label="GitHub"
        />
      </template>
    </UFooter>
  </div>
</template>

<script setup lang="ts">
const colorMode = useColorMode()

const navItems = [
  {label: 'Home', to: '/'},
  {label: 'Documentation', to: '/docs'},
]
</script>
