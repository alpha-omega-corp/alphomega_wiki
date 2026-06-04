# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository layout

This git repo (`pages/`) holds two independently-run applications that talk to each other over HTTP:

- **`docs-app/`** — a Nuxt 4 frontend ("DevDocs") plus a thin Nitro server layer that proxies to the Go backend. This is the primary working directory.
- **`server/`** — a Go backend (`module wiki`) that owns persistence: Postgres via the [bun](https://bun.uptrace.dev/) ORM, routed with `bunrouter`.

The root `docker-compose.yml`, `package.json`, `src/`, and `tsconfig.json` are stale/unrelated scaffolding — the real Postgres compose file is `server/docker-compose.yaml`.

## Running the stack

Three processes, in order:

```bash
# 1. Postgres (db=user=password = "wiki", exposed on localhost:5432)
cd server && docker compose -f docker-compose.yaml up -d

# 2. Go backend — serves http://localhost:1000
cd server && go run ./cmd

# 3. Nuxt frontend — serves http://localhost:3000
cd docs-app && npm run dev
```

There are **no tests** and no linter configured in either project.

Notes:
- `docs-app`'s `npm run build` is just `tsc` (typecheck), **not** `nuxt build` — it does not produce a deployable Nuxt output.
- **The Go server wipes the database on every startup.** `cmd/main.go` calls `db.ResetModel` (drops + recreates the `pages` table) and reloads seed rows from `cmd/migration.yml`. Any data created through the UI is lost on backend restart, and `migration.yml` is the source of truth for seed pages.

## Architecture: request flow

```
Browser → Nuxt pages/components
        → Nitro routes in docs-app/server/api/*   (BFF / proxy layer)
        → Go backend at http://localhost:1000     (hardcoded in each Nitro handler)
        → Postgres
```

The Nitro layer (`docs-app/server/api/`) is mostly a pass-through to the Go API, with two responsibilities of its own:

- **`page/[id].ts`** is the only non-trivial proxy: it fetches the raw page from Go, then parses the markdown `content` into a renderable `body` + `toc` using `@nuxtjs/mdc/runtime` with a server-side Shiki highlighter (`golang` is aliased to `go`). The frontend renders this `body` via `<ContentRenderer>`.
- **`page/[id].put.ts`** merges the partial client body over the current upstream page (`{...current, ...body}`) before forwarding, so callers can PATCH a single field (e.g. just `content` or just `title`).

Backend endpoints (registered in `server/cmd/main.go`): `GET /pages`, `POST /page`, `GET /page/:id`, `PUT /page/:id`, `DELETE /page/:id`.

## Architecture: slug-vs-id duality (important)

Pages are **keyed by `id`** in the API/database, but **URLs are slug-based** (`/docs/parent/child`) for readability. The frontend bridges this:

- The full page tree is fetched **once** via `useAsyncData('pages', ...)`. This `'pages'` cache key is shared between `app.vue` (for the search palette) and `PageList.vue` (the sidebar) so the tree is only loaded once.
- `pages/docs/[...slug].vue` resolves a slug path → `id` by walking that cached tree (`findIdByPath`), then fetches the page by id.
- After any mutation (create/rename/delete/save), components call `refreshNuxtData('pages')` to rebuild the shared tree.

When adding features that touch pages, follow this pattern: mutate by id, navigate by slug, and refresh the `'pages'` key.

## Architecture: the page tree (Go side)

`models.Page` is self-referential via `ParentID *int64` (`rel:has-many` on `id=parent_id`). Because bun relations only nest to a fixed depth, `handlers/page.go::GetAll` deliberately **does not** use relation loading — it fetches every page flat (excluding `content`/`excerpt` for the list) and assembles the unbounded parent/child tree in memory. Pages with no known parent become roots.

Persistence goes through a generic `utils.Repository[E]` (in `server/utils/repository.go`) — `GetOne`, `GetOneBy`, `GetAll(tweak, relations...)`, `Create`, `Update`, `Delete`. `utils/http.go` wraps handlers with `Response`/`JSON`/`Error` helpers and extracts the `:id` path param via `GetParams`. Reuse these rather than writing raw bun queries in controllers.

## Architecture: the editor

`docs-app/components/editor/Editor.vue` is a TipTap editor built on `@nuxt/ui`'s `UEditor`, opened from a drawer in `Page.vue` (toggle with the `o` shortcut). It works in `content-type="markdown"`, so the editor's value round-trips as markdown — the same string stored in the Go `content` column. Saving PUTs `{ content }` to `/api/page/:id` then emits `saved`, which triggers a `refreshNuxtData('docs-<slug>')` to re-render.

Custom TipTap pieces live alongside it:
- `EditorImageUploadExtension.ts` — a drag/click upload placeholder node that inlines images as base64 data URLs (no upload backend).
- `EditorUseCompletion.ts` — AI text actions (fix/extend/translate/…) that POST to **`/api/ai/complete`**. ⚠️ This endpoint does **not exist** in `server/api/` yet; the handler silently swallows the error, so these toolbar actions are currently no-ops until that route is implemented.

## Authentication (OAuth via AOAuth)

Login is handled by the **AOAuth** OAuth 2.0 / OIDC server using the
Authorization Code + PKCE flow. It is implemented entirely in the Nitro layer —
tokens are sealed in an encrypted, httpOnly session cookie (h3 `useSession`) and
**never reach the browser**; only `{ user }` is exposed via `/api/auth/session`.

Flow / files:
- `server/api/auth/login.get.ts` → generates PKCE + state, redirects to AOAuth `/authorize`.
- `server/routes/callback.get.ts` → validates state, exchanges the code for tokens
  (server-side, with the client secret), fetches `/userinfo`, seals the session.
- `server/api/auth/{session.get,logout.post}.ts` → read / clear the session.
- `server/utils/session.ts` → `authSession()` + `oauthConfig()` (reads `AOAUTH_*`
  from `process.env` at runtime).
- `composables/useAuth.ts` + `components/AuthButton.vue` → header login/logout UI.
- `middleware/auth.global.ts` → gates the whole app when `AUTH_REQUIRED=true` (off by default).

Config is via env (see `docs-app/.env.example`); defaults match AOAuth's seeded
`demo-client` (redirect `http://localhost:3000/callback`), so with AOAuth on
`:8080` and this app on `:3000` it works with no setup. Run AOAuth from its repo
(`AOAUTH_ADMIN_API_KEY=… go run ./cmd/aoauth`), then sign in with the seeded
`alice` / `password`.

**Defense in depth (Go backend).** Mutations are enforced in two places: the
Nitro proxy requires a session and forwards the access token as a Bearer header
(`server/utils/session.ts::requireAuth`), and the **wiki backend independently
verifies** that token. `server/app/middleware/auth.go` (`JWTVerifier`) validates
the RS256 JWT against AOAuth's JWKS (`${AOAUTH_ISSUER}/.well-known/jwks.json`,
default `:8080`), checking signature / `exp` / `iss`. It wraps the
`POST/PUT/DELETE /page` routes in `cmd/main.go`; `GET` routes stay public. The
verifier fails closed (401 if AOAuth is unreachable) and refetches the JWKS on
an unknown `kid` to handle key rotation. Note: the backend resets its DB on
startup, so **restart it** to pick up the new middleware.

## Conventions

- Frontend is TypeScript + `<script setup>`; the shared `Page` type lives in `docs-app/types/page.ts` (extends the API shape with optional `body`/`toc` present only on the single-page response).
- UI is `@nuxt/ui` v4 components (`U*`) with Tailwind utility classes; dark mode via `useColorMode()`.
- The Nitro proxy target (`http://localhost:1000`) is hardcoded in each handler — change it in all of `docs-app/server/api/**` if the backend moves.

## Guardrails

Quick do/don't rules (data flow, gotchas, UI) are kept in a dedicated file:

@.claude/rules/project-guardrails.md