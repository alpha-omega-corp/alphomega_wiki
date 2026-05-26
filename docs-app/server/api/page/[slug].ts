import { parseMarkdown, createShikiHighlighter } from '@nuxtjs/mdc/runtime'
import rehypeSlug from 'rehype-slug'
import { bundledLanguages, bundledThemes } from 'shiki'
import type { Page } from '~/types/page'

const highlighter = createShikiHighlighter({
  bundledLangs: bundledLanguages,
  bundledThemes: bundledThemes,
})

export default defineEventHandler(async (event) => {
  const slug = getRouterParam(event, 'slug')

  const page = await $fetch<Page>(`http://localhost:1000/page/${slug}`)

  const parsed = await parseMarkdown(page.content, {
    highlight: {
      highlighter,
      theme: {
        default: 'github-light',
        dark: 'github-dark',
      },
    },
    rehype: {
      plugins: {
        'rehype-slug': { instance: rehypeSlug },
      },
    },
  })

  // Extract shiki CSS from the <style> node injected by rehypeHighlight, then remove it from body
  // so it can be injected into <head> by the page component instead (fixes CSR navigation highlighting)
  const styleIndex = parsed.body?.children?.findIndex((child: any) => child.tag === 'style') ?? -1
  let highlightCSS = ''
  if (styleIndex !== -1) {
    highlightCSS = (parsed.body.children[styleIndex] as any)?.children?.[0]?.value ?? ''
    parsed.body.children.splice(styleIndex, 1)
  }

  return {
    title: page.title,
    slug: page.slug,
    tags: page.tags,
    content: page.content,
    body: parsed.body,
    toc: parsed.toc,
    highlightCSS,
  }
})