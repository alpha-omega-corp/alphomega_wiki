import type {Page} from '~/types/page'

export default defineEventHandler(async (event) => {
    const slug = getRouterParam(event, 'slug')
    const body = await readBody<Partial<Page>>(event)

    const current = await $fetch<Page>(`http://localhost:1000/page/${slug}`)

    return $fetch<Page>(`http://localhost:1000/page/${slug}`, {
        method: 'PUT',
        body: {...current, ...body},
    })
})