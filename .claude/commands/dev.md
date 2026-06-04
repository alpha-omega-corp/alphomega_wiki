---
description: Start the full local stack (Postgres, Go backend, Nuxt frontend)
allowed-tools: Bash
---

Bring up the full local stack, in order, each in the background:

1. **Postgres** — `cd server && docker compose -f docker-compose.yaml up -d`
   (db=user=password = `wiki`, exposed on `localhost:5432`)
2. **Go backend** — `cd server && go run ./cmd` (serves `http://localhost:1000`).
   ⚠️ This **drops and re-seeds** the `pages` table from `cmd/migration.yml` on every start.
3. **Nuxt frontend** — `cd docs-app && npm run dev` (serves `http://localhost:3000`).

Run the backend and frontend as background processes so they keep running. After
starting, confirm each is up (curl `http://localhost:1000/pages` and report the
frontend URL). Do not restart the backend casually — it wipes any data created via the UI.
