// POST /api/favorite — star a page for the current user. Body: { page_id }.
// Requires an authenticated session; the access token is forwarded so the
// backend records the favorite against the verified user (never the body).
export default defineEventHandler(async (event) => {
  const { accessToken } = await requireAuth(event)
  const body = await readBody<{ page_id: number }>(event)
  return $fetch('http://localhost:1000/favorite', {
    method: 'POST',
    body,
    headers: bearer(accessToken),
  })
})