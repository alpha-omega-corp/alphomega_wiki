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


    // The app is served by the Nitro node server (SSR), not statically
    // prerendered. @nuxt/content auto-registers a prerender route for its SQL
    // dump; with no local `content/` files that route is empty (204), which would
    // otherwise abort `nuxt build`. Don't fail the build on it.
    nitro: {
        prerender: {
            failOnError: false,
        },
    },

    devtools: {enabled: true},

    // OAuth login against the AOAuth server. Secrets stay server-side; only
    // `public` is exposed to the browser.
    runtimeConfig: {
        oauth: {
            issuer: process.env.AOAUTH_ISSUER || 'http://localhost:8080',
            clientId: process.env.AOAUTH_CLIENT_ID || 'demo-client',
            clientSecret: process.env.AOAUTH_CLIENT_SECRET || 'demo-secret',
            redirectUri: process.env.AOAUTH_REDIRECT_URI || 'http://localhost:3000/callback',
        },
        session: {
            // Seals the session cookie. MUST be >= 32 chars; override in production.
            password: process.env.AUTH_SESSION_PASSWORD || 'dev-only-insecure-session-password-change-me!!',
        },
        public: {
            // When true, every page requires authentication (middleware/auth.global.ts).
            authRequired: process.env.AUTH_REQUIRED === 'true',
        },
    },

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