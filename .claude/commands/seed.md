---
description: Inspect or edit the seed page data
allowed-tools: Read, Edit, Bash
---

Seed pages live in `server/cmd/migration.yml` and are the source of truth for what
exists after a backend restart (the Go server wipes + reloads them on every start).

Each entry is a `Page` row: `id`, `slug`, `title`, `published`, optional `parent_id`
(for nesting in the sidebar tree), `tags`, and markdown `content`.

Open `server/cmd/migration.yml` and help with the requested change: $ARGUMENTS

When adding a page, pick an unused `id`, set a unique `slug`, and use `parent_id` to
nest it under an existing page. Restart the Go backend to load the new fixture.