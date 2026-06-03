<template>
  <div class="min-h-screen flex flex-col bg-gray-50 dark:bg-gray-950">
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
    <main class="flex-1">
      <slot/>
    </main>


    <!-- Footer -->
    <UFooter>
      <template #left>
        <p class="text-muted text-sm">Copyright © {{ new Date().getFullYear() }}</p>
      </template>

      <UNavigationMenu :items="items" variant="link"/>

      <template #right>
        <UButton
            icon="i-simple-icons-discord"
            color="neutral"
            variant="ghost"
            to="https://go.nuxt.com/discord"
            target="_blank"
            aria-label="Discord"
        />
        <UButton
            icon="i-simple-icons-x"
            color="neutral"
            variant="ghost"
            to="https://go.nuxt.com/x"
            target="_blank"
            aria-label="X"
        />
        <UButton
            icon="i-simple-icons-github"
            color="neutral"
            variant="ghost"
            to="https://github.com/nuxt/nuxt"
            target="_blank"
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

import type {NavigationMenuItem} from '@nuxt/ui'

const items: NavigationMenuItem[] = [
  {
    label: 'Figma Kit',
    to: 'https://go.nuxt.com/figma-ui',
    target: '_blank'
  },
  {
    label: 'Playground',
    to: 'https://stackblitz.com/edit/nuxt-ui',
    target: '_blank'
  },
  {
    label: 'Releases',
    to: 'https://github.com/nuxt/ui/releases',
    target: '_blank'
  }
]
</script>
