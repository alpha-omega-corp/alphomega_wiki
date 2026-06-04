import type {MDCParserResult} from "@nuxtjs/mdc";

export type PageTocLinks = NonNullable<MDCParserResult['toc']>['links']

export interface Page {
    id: number
    slug: string
    title: string
    icon?: string
    published: boolean
    tags?: string[]
    content: string
    parent_id?: number | null
    position?: number
    pages?: Page[]
    // Present on the single-page API response (`/api/page/:slug`), where the
    // markdown is parsed into a renderable body + table of contents.
    body?: MDCParserResult['body']
    toc?: MDCParserResult['toc']
}