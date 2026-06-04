import {mergeAttributes, Node, parseAttributes} from '@tiptap/core'
import {VueNodeViewRenderer} from '@tiptap/vue-3'
import MdcCardView from './MdcCardView.vue'

// A TipTap node that round-trips Nuxt Content's MDC block-component syntax for
// `::card{...}` and renders it live (via the real content/WCard.vue component)
// inside the editor instead of showing the raw `::card` text.
//
//   ::card{icon="i-lucide-triangle-alert" variant="warning"}
//   Body **markdown** here.
//   ::
//
// The markdown plumbing (markdownTokenizer / parseMarkdown / renderMarkdown) is
// consumed by @tiptap/markdown, which UEditor wires up for content-type="markdown".

const CARD_ATTRS = ['icon', 'variant'] as const

const DEFAULT_ATTRS = {icon: '', variant: 'default'}

export const MdcCard = Node.create({
    name: 'mdcCard',
    group: 'block',
    content: 'block+',
    defining: true,
    isolating: true,
    draggable: true,

    addAttributes() {
        return {
            icon: {default: ''},
            variant: {default: 'default'},
        }
    },

    parseHTML() {
        return [{tag: 'div[data-mdc="card"]'}]
    },

    renderHTML({HTMLAttributes}) {
        return ['div', mergeAttributes(HTMLAttributes, {'data-mdc': 'card'}), 0]
    },

    addNodeView() {
        return VueNodeViewRenderer(MdcCardView)
    },

    // --- Markdown round-trip -------------------------------------------------
    markdownTokenName: 'mdcCard',

    markdownTokenizer: {
        name: 'mdcCard',
        level: 'block',
        // Hint marked where a card could start so it doesn't swallow it into a paragraph.
        start: (src: string) => src.indexOf('::card'),
        tokenize(src, _tokens, lexer) {
            // ::card{attrs}\n  <body...>  \n::
            const match = /^::card[ \t]*(?:\{([^}]*)\})?[^\n]*\n([\s\S]*?)\n::[ \t]*(?:\n|$)/.exec(src)
            if (!match) return undefined

            const attributes = {...DEFAULT_ATTRS, ...parseAttributes(match[1] || '')}
            const body = (match[2] || '').replace(/^\n+/, '').replace(/\n+$/, '')

            return {
                type: 'mdcCard',
                raw: match[0],
                attributes,
                // Parse the body as block-level markdown so it becomes editable content.
                tokens: body ? lexer.blockTokens(body) : [],
            }
        },
    },

    parseMarkdown(token, helpers) {
        const attrs = {...DEFAULT_ATTRS, ...(token.attributes || {})}
        const content = helpers.parseChildren(token.tokens || [])
        return {
            type: 'mdcCard',
            attrs,
            // `block+` needs at least one child; keep an empty paragraph for empty cards.
            content: content.length ? content : [{type: 'paragraph'}],
        }
    },

    renderMarkdown(node, helpers) {
        const attrs = node.attrs || {}
        const parts: string[] = []
        for (const key of CARD_ATTRS) {
            const value = attrs[key]
            // Skip empty values and the default variant to keep the markdown tidy.
            if (!value || (key === 'variant' && value === 'default')) continue
            parts.push(`${key}="${String(value).replace(/"/g, '\\"')}"`)
        }
        const attrString = parts.length ? `{${parts.join(' ')}}` : ''
        const inner = helpers.renderChildren(node).trim()
        return `::card${attrString}\n${inner}\n::`
    },
})
