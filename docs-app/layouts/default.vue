<template>
  <div class="h-dvh flex flex-col overflow-hidden bg-gray-50 dark:bg-gray-950">
    <!-- Top bar: brand + auth control, visible on every page. -->
    <header
        class="flex-shrink-0 h-14 flex items-center justify-between gap-4 px-4 sm:px-6 border-b border-default bg-white/80 dark:bg-gray-900/80 backdrop-blur"
    >
      <NuxtLink to="/" class="flex items-center gap-2 font-semibold">
        <UIcon name="i-lucide-book-open-text" class="size-5 text-primary-500"/>
        <span>DevDocs</span>
      </NuxtLink>

      <div class="flex items-center gap-1.5">
        <UButton
            to="/"
            icon="i-lucide-house"
            label="Home"
            color="neutral"
            variant="ghost"
            size="sm"
        />
        <UButton
            :icon="isDark ? 'i-lucide-moon' : 'i-lucide-sun'"
            color="neutral"
            variant="ghost"
            size="sm"
            aria-label="Toggle theme"
            @click="isDark = !isDark"
        />
        <AuthButton/>
      </div>
    </header>

    <!-- App Content (the scroll boundary; individual columns scroll inside it) -->
    <main class="flex-1 min-h-0 overflow-y-auto">
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
import type {NavigationMenuItem} from '@nuxt/ui'

const colorMode = useColorMode()
const isDark = computed({
  get: () => colorMode.value === 'dark',
  set: (v) => {
    colorMode.preference = v ? 'dark' : 'light'
  },
})

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
