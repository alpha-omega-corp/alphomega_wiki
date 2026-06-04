import type { Page } from '~/types/page'

// A favorite row as returned by the Go backend (GET /api/favorites). Scoped to
// the calling user server-side via their verified access token.
export interface Favorite {
  id: number
  user_sub: string
  page_id: number
  created_at: string
}

// loadFavorites performs the single, shared fetch of the authenticated user's
// favorites under the 'favorites' key. Call it exactly once from the component
// that owns the page tree (PageList) — never per row, or the watch on `loggedIn`
// would refire the request once per PageItem. Logged-out users get an empty
// list, so the request never 401s.
export function loadFavorites() {
  const { loggedIn } = useAuth()
  return useAsyncData<Favorite[]>(
    'favorites',
    () => (loggedIn.value ? $fetch<Favorite[]>('/api/favorites') : Promise.resolve([])),
    { default: () => [], server: false, watch: [loggedIn] },
  )
}

// useFavorites reads the shared 'favorites' list (without triggering its own
// request — that's loadFavorites' job) and exposes isFavorite / toggle. Safe to
// call from every PageItem. The user is always resolved from the session
// server-side; it is never passed from here.
export function useFavorites() {
  const toast = useToast()

  // useNuxtData reads the cached ref filled by loadFavorites; it reactively
  // updates after every refreshNuxtData('favorites') without re-fetching here.
  const { data } = useNuxtData<Favorite[]>('favorites')

  // Ordered as returned by the backend (newest first).
  const favorites = computed(() => data.value ?? [])
  const favoriteIds = computed(() => new Set(favorites.value.map(f => f.page_id)))
  const isFavorite = (pageId: number) => favoriteIds.value.has(pageId)

  // Star/unstar a page for the current user, then re-sync the shared list.
  async function toggle(page: Page) {
    const wasFavorite = isFavorite(page.id)
    try {
      if (wasFavorite) {
        await $fetch(`/api/favorite/${page.id}`, { method: 'DELETE' })
      } else {
        await $fetch('/api/favorite', { method: 'POST', body: { page_id: page.id } })
      }
      await refreshNuxtData('favorites')
      toast.add({
        title: wasFavorite ? 'Removed from favorites' : 'Added to favorites',
        description: page.title,
        color: 'success',
        icon: wasFavorite ? 'i-lucide-star-off' : 'i-lucide-star',
      })
    } catch (e: any) {
      toast.add({
        title: 'Could not update favorites',
        description: e?.data?.message ?? e?.message ?? 'Failed to update favorites.',
        color: 'error',
        icon: 'i-lucide-triangle-alert',
      })
    }
  }

  return { favorites, favoriteIds, isFavorite, toggle }
}