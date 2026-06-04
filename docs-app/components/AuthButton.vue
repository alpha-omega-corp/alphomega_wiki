<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui'

const { user, loggedIn, login, logout } = useAuth()
const route = useRoute()

const items = computed<DropdownMenuItem[][]>(() => [
  [
    {
      label: user.value?.name || user.value?.email || 'Signed in',
      type: 'label',
    },
  ],
  [
    {
      label: 'Sign out',
      icon: 'i-lucide-log-out',
      onSelect: () => logout(),
    },
  ],
])
</script>

<template>
  <UButton
    v-if="!loggedIn"
    icon="i-lucide-log-in"
    label="Sign in"
    color="primary"
    size="sm"
    @click="login(route.fullPath)"
  />

  <UDropdownMenu v-else :items="items" :content="{ align: 'end' }">
    <UButton color="neutral" variant="ghost" size="sm" trailing-icon="i-lucide-chevron-down">
      <UAvatar icon="i-lucide-user" size="2xs" />
      <span class="hidden sm:inline max-w-40 truncate">{{ user?.name || user?.email }}</span>
    </UButton>
  </UDropdownMenu>
</template>
