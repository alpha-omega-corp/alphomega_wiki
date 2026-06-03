import {createShikiHighlighter, parseMarkdown} from '@nuxtjs/mdc/runtime'
import {bundledLanguages, bundledThemes} from 'shiki'
import type {Page} from '~/types/page'

const highlighter = createShikiHighlighter({
    bundledLangs: {
        ...bundledLanguages,
        // Register the `golang` alias so ```golang fences highlight like ```go.
        golang: bundledLanguages.go,
    },
    bundledThemes: bundledThemes,
})

export default defineEventHandler(async (event) => {
    const slug = getRouterParam(event, 'slug')

    const page: Page = await $fetch<Page>(`http://localhost:1000/page/${slug}`)

    const parsed = await parseMarkdown(page.content, {
        highlight: {
            highlighter,
            theme: {
                default: 'github-light',
                dark: 'github-dark',
            },
        },
    })

    return {
        title: page.title,
        slug: page.slug,
        tags: page.tags,
        content: page.content,
        body: parsed.body,
        toc: parsed.toc,
    }
})