// DELETE /api/page/:id — delete a page. Requires an authenticated session; the
// AOAuth access token is forwarded to the backend.
export default defineEventHandler(async (event) => {
    const { accessToken } = await requireAuth(event)
    const id = getRouterParam(event, 'id')
    return $fetch(`http://localhost:1000/page/${id}`, {
        method: 'DELETE',
        headers: bearer(accessToken),
    })
})
