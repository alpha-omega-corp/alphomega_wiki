export interface NavPage {
  id: number
  slug: string
  title: string
  tags?: string[]
  published?: boolean
  pages?: NavPage[]
}