import type {Page} from '~/types/page'

export default defineEventHandler(() => {
    return $fetch<Page[]>(`${backendUrl()}/pages`)
})
