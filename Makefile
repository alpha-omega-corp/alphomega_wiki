IMAGE   := wiki

# Root of the monorepo (one level up from server/).
ROOT := $(abspath ..)

.PHONY: build run

# Build context is the repo root so the Dockerfile can reach both server/ and docs-app/.
build:
	docker build -t $(IMAGE) .

# Run the combined backend+frontend container.
run: build
	docker compose up -d
