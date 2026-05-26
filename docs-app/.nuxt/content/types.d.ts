import type { PageCollectionItemBase, DataCollectionItemBase } from '@nuxt/content'

declare module '@nuxt/content' {
   interface DocsCollectionItem extends PageCollectionItemBase {
    title: string
    description?: string
  }
  

  interface PageCollections {
    docs: DocsCollectionItem
  }

  interface Collections {
    docs: DocsCollectionItem
  }
}
