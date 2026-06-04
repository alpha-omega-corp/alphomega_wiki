// Optional app-wide gate. Disabled unless AUTH_REQUIRED=true. When enabled,
// unauthenticated visitors are sent through the OAuth login flow.
export default defineNuxtRouteMiddleware(async (to) => {
  const config = useRuntimeConfig()
  if (!config.public.authRequired) return

  // Never gate the auth endpoints or the OAuth callback themselves.
  if (to.path.startsWith('/api/') || to.path === '/callback') return

  const headers = import.meta.server ? useRequestHeaders(['cookie']) : undefined
  const session = await $fetch<{ user: unknown }>('/api/auth/session', { headers }).catch(
    () => ({ user: null }),
  )

  if (!session.user) {
    return navigateTo(`/api/auth/login?returnTo=${encodeURIComponent(to.fullPath)}`, {
      external: true,
    })
  }
})
