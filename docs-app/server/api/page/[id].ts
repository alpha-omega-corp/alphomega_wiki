import {createShikiHighlighter, parseMarkdown} from '@nuxtjs/mdc/runtime'
import {bundledLanguages, bundledThemes} from 'shiki'
import type {Page} from '~/types/page'

const highlighter = createShikiHighlighter({
    bundledLangs: {
        ...bundledLanguages,
        golang: bundledLanguages.go,
    },
    bundledThemes: bundledThemes,
})

export default defineEventHandler(async (event) => {
    const id = getRouterParam(event, 'id')

    const page: Page = await $fetch<Page>(`http://localhost:1000/page/${id}`)

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