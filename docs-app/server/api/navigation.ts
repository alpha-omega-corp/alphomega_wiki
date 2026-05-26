import type { NavPage } from '~/types/navigation'

export default defineCachedEventHandler(async () => {
  return $fetch<NavPage[]>('http://localhost:1000/navigation')
}, {
  getKey: () => 'navigation',
})