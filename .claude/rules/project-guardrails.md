# Project guardrails

Crisp do/don't rules for this repo. Background/architecture is in the root `CLAUDE.md`.

## Data & API
- Mutate pages **by `id`**, navigate **by `slug`**. To go from a slug path to an id,
  walk the cached `'pages'` tree (`findIdByPath`) — never assume a 1:1 slug↔id mapping.
- After any page mutation (create / rename / delete / save), call
  `refreshNuxtData('pages')` so the shared sidebar/search tree rebuilds. For the
  currently-viewed page also refresh `docs-<slug>`.
- The Nitro routes in `docs-app/server/api/**` are a proxy to the Go backend at the
  hardcoded `http://localhost:1000`. New page operations need **both** a Nitro route
  and a Go endpoint — don't add DB logic in Nitro.
- On the Go side, go through `utils.Repository[E]` and the `utils.Response/JSON/Error`
  helpers; don't write raw bun queries or hand-roll JSON in controllers.

## Gotchas
- The Go server **wipes and re-seeds the DB on every startup** (`ResetModel` +
  `cmd/migration.yml`). Don't restart it expecting UI-created data to survive; edit
  `migration.yml` for durable seed data.
- `docs-app`'s `npm run build` is just `tsc` (typecheck), not a Nuxt production build.
- The editor's AI actions POST to `/api/ai/complete`, which **does not exist yet** —
  failures are swallowed silently. Implement that route before relying on them.
- Page `content` is markdown end-to-end: the editor round-trips markdown, the Go
  `content` column stores markdown, and `page/[id].ts` parses it to `body`/`toc` for
  rendering. Keep it markdown — don't store HTML.

## UI
- Use `@nuxt/ui` v4 `U*` components and Tailwind utilities; dark mode via `useColorMode()`.
- Main docs content width is governed by the `UContainer` in `pages/docs/[...slug].vue`,
  not by the `Page` component itself.