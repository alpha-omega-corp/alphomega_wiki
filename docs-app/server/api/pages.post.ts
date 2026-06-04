// POST /api/pages — create a page. Requires an authenticated session; the
// AOAuth access token is forwarded to the backend.
export default defineEventHandler(async (event) => {
  const { accessToken } = await requireAuth(event)
  const body = await readBody(event)
  return $fetch('http://localhost:1000/page', {
    method: 'POST',
    body,
    headers: bearer(accessToken),
  })
})
