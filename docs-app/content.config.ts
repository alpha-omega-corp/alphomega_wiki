import { defineContentConfig, defineCollection, z } from '@nuxt/content'

export default defineContentConfig({
  collections: {
    docs: defineCollection({
      type: 'page',
      schema: z.object({
        title: z.string(),
        description: z.string().optional(),
      }),
    }),
  },
})
