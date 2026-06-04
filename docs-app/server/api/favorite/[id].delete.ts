// DELETE /api/favorite/:id — unstar a page (by page id) for the current user.
// Requires an authenticated session; the access token is forwarded so the
// backend removes the favorite for the verified user (never from the body).
export default defineEventHandler(async (event) => {
  const { accessToken } = await requireAuth(event)
  const id = getRouterParam(event, 'id')
  return $fetch(`http://localhost:1000/favorite/${id}`, {
    method: 'DELETE',
    headers: bearer(accessToken),
  })
})