import type {Page} from '~/types/page'
// PUT /api/page/:id — update a page (save / rename / reorder). Requires an
// authenticated session; the AOAuth access token is forwarded to the backend.
export default defineEventHandler(async (event) => {
    const {accessToken} = await requireAuth(event)

    // The upstream page routes are keyed by id; the client passes the page id as
    // the path segment.
    const id = getRouterParam(event, 'id')
    const body = await readBody<Partial<Page>>(event)

    const current = await $fetch<Page>(`${backendUrl()}/page/${id}`)

    return $fetch<Page>(`${backendUrl()}/page/${id}`, {
        method: 'PUT',
        body: {...current, ...body},
        headers: bearer(accessToken),
    })
})
