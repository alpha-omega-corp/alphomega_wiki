// GET /api/favorites — the current user's favorited pages. Requires an
// authenticated session; the AOAuth access token is forwarded so the backend
// can scope the list to the calling user.
export default defineEventHandler(async (event) => {
  const { accessToken } = await requireAuth(event)
  return $fetch('http://localhost:1000/favorites', {
    headers: bearer(accessToken),
  })
})