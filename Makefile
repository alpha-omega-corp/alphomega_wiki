IMAGE   := wiki

# Root of the monorepo (one level up from server/).
ROOT := $(abspath ..)

.PHONY: build run


# Build context is the repo root so the Dockerfile can reach both server/ and docs-app/.
build:
	docker build -t $(IMAGE) .

# Run the combined backend+frontend container.
run:
	docker compose up -d --force-recreate

tag:
	docker tag wiki:latest ghcr.io/nanstis/wiki:latest

push: build tag
	docker push ghcr.io/nanstis/wiki:latest

prod:
	sudo --preserve-env=AOAUTH_ADMIN_API_KEY docker compose up -d
