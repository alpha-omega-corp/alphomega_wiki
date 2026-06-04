interface SessionUser {
  sub: string
  name?: string
  email?: string
  email_verified?: boolean
}

// useAuth exposes the current login state and login/logout actions. The session
// is fetched once (shared cache key) and re-used across components; SSR forwards
// the request cookie so the state is correct on first render.
export function useAuth() {
  const headers = import.meta.server ? useRequestHeaders(['cookie']) : undefined

  const { data, refresh } = useFetch<{ user: SessionUser | null }>('/api/auth/session', {
    key: 'auth-session',
    headers,
    default: () => ({ user: null }),
  })

  const user = computed(() => data.value?.user ?? null)
  const loggedIn = computed(() => !!user.value)

  function login(returnTo?: string) {
    const target = returnTo && returnTo.startsWith('/') ? returnTo : '/'
    const url = `/api/auth/login?returnTo=${encodeURIComponent(target)}`
    if (import.meta.client) window.location.assign(url)
  }

  async function logout() {
    await $fetch('/api/auth/logout', { method: 'POST' })
    await refresh()
    await navigateTo('/')
  }

  return { user, loggedIn, login, logout, refresh }
}
