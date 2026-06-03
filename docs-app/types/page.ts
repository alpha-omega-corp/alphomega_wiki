import type {MDCParserResult} from "@nuxtjs/mdc";

export interface Page {
    id: number
    slug: string
    title: string
    published: boolean
    tags?: string[]
    content: string
    pages?: Page[]
    parsedContent: MDCParserResult
}