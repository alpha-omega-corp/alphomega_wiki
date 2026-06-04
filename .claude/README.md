# .claude/

Claude Code configuration for this repo. Loaded automatically when working anywhere
under `pages/` (the git root).

- **`settings.json`** — committed, team-wide: permission allowlist for common dev
  commands (npm/go/docker/read-only git) plus a small deny list. Edit via the
  `/config` UI or by hand.
- **`settings.local.json`** — personal, git-ignored overrides for this project.
- **`commands/`** — project slash commands:
  - `/dev` — start the full local stack (Postgres → Go backend → Nuxt).
  - `/typecheck` — `tsc` the frontend and `go build`/`go vet` the backend.
  - `/seed` — inspect/edit the seed pages in `server/cmd/migration.yml`.
- **`rules/project-guardrails.md`** — crisp do/don't rules, imported by the root
  `CLAUDE.md`. The root `CLAUDE.md` holds the architecture overview.