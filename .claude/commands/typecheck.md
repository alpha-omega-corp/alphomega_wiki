---
description: Typecheck the frontend and build-check the Go backend
allowed-tools: Bash
---

Run static checks across both apps and report failures concisely:

- **Frontend** (`docs-app/`): `cd docs-app && npm run build` (runs `tsc`).
- **Backend** (`server/`): `cd server && go build ./... && go vet ./...`.

There is no test suite in either project, so these are the primary correctness gates.