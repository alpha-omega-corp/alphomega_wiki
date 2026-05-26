import type { Page } from '~/types/page'

export default defineEventHandler(() => {
  return $fetch<Page[]>('http://localhost:1000/navigation')
})
