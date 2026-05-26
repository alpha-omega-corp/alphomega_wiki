export default defineNuxtConfig({
  srcDir: '.',

  modules: [
    '@nuxt/ui',
    '@nuxt/content',
  ],

  vite: {
      optimizeDeps: {
          include: [
              '@nuxt/ui > prosemirror-gapcursor',
              '@nuxt/ui > prosemirror-model',
              '@nuxt/ui > prosemirror-state',
              '@nuxt/ui > prosemirror-transform',
              '@nuxt/ui > prosemirror-view',
              '@tiptap/core',
              '@tiptap/extension-emoji',
              '@tiptap/extension-text-align',
              '@tiptap/vue-3',
              '@vue/devtools-core',
              '@vue/devtools-kit',
              'tiptap-extension-code-block-shiki'
          ]
      }
  },


  devtools: { enabled: true },

  content: {
    experimental: {
      sqliteConnector: 'native',
    },
    build: {
      markdown: {
        highlight: {
          theme: {
            default: 'github-light',
            dark: 'github-dark',
          },
          langs: [
            'javascript', 'typescript', 'vue', 'jsx', 'tsx',
            'html', 'css', 'json', 'bash', 'markdown', 'yaml', 'python',
          ],
        },
      },
    },
  },

  css: ['~/assets/css/main.css'],

  compatibilityDate: '2025-05-19',
})