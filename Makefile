.PHONY: help build-backend build-frontend sync-frontend build start dev test clean docker
# scaffold-only:start
.PHONY: create
# scaffold-only:end

SHELL := /bin/bash

VERSION ?= latest
IMG_NAME ?= go_web_template
COMMIT ?= $(shell git rev-parse --short=8 HEAD 2>/dev/null || echo unknown)
BUILDTIME ?= $(shell date '+%Y-%m-%d %H:%M:%S %z')
LDFLAGS := -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X 'main.BuildTime=$(BUILDTIME)'

BACKEND_OUT ?= bin/go_web_template
CREATE_OUT ?= bin/create-go-web
WEB_DIR := web
WEB_DIST := $(WEB_DIR)/dist
WEB_ASSETS_DIR := assets/web/v1

help:
	@echo "Targets:"
	@echo "  build-backend     Build Go binary -> $(BACKEND_OUT)"
	@echo "  build-frontend    Build Vue frontend -> $(WEB_DIST)"
	@echo "  sync-frontend     Copy $(WEB_DIST) -> $(WEB_ASSETS_DIR)"
	@echo "  build             Build backend+frontend and sync assets"
	@echo "  start             Build+sync then run (http://127.0.0.1:1423/v1/web/)"
	@echo "  dev               Run server (go run) + Vite dev server"
	@echo "  test              Run Go tests"
	@echo "  clean             Remove build artifacts"
	@echo "  docker            Build image $(IMG_NAME):$(VERSION) / latest / $(COMMIT)"
# scaffold-only:start
	@echo "  create            Build $(CREATE_OUT) and scaffold a new project"
# scaffold-only:end

# scaffold-only:start
create:
	@mkdir -p "$$(dirname "$(CREATE_OUT)")"
	go build -o "$(CREATE_OUT)" ./cmd/create
	./"$(CREATE_OUT)"
# scaffold-only:end

build-backend:
	@set -euo pipefail; \
	mkdir -p "$$(dirname "$(BACKEND_OUT)")"; \
	go build -ldflags "$(LDFLAGS)" -o "$(BACKEND_OUT)" .

build-frontend:
	@set -euo pipefail; \
	cd "$(WEB_DIR)"; \
	npm install; \
	npm run build

sync-frontend:
	@set -euo pipefail; \
	if [ ! -d "$(WEB_DIST)" ]; then \
		echo "Missing $(WEB_DIST). Run: make build-frontend"; \
		exit 1; \
	fi; \
	mkdir -p "$(WEB_ASSETS_DIR)"; \
	rm -rf "$(WEB_ASSETS_DIR)"/*; \
	cp -R "$(WEB_DIST)"/. "$(WEB_ASSETS_DIR)"/

build: build-backend build-frontend sync-frontend

start: build
	@set -euo pipefail; \
	if [ -f .env ]; then set -a && . ./.env && set +a; fi; \
	./"$(BACKEND_OUT)"

dev:
	@set -euo pipefail; \
	if [ -f .env ]; then set -a && . ./.env && set +a; fi; \
	( cd "$(WEB_DIR)" && npm install && npm run dev ) & \
	FRONT_PID="$$!"; \
	trap 'kill "$$FRONT_PID" 2>/dev/null || true' EXIT; \
	go run -ldflags "$(LDFLAGS)" .; \
	wait "$$FRONT_PID"

test:
	go test ./...

clean:
	@rm -rf "$(BACKEND_OUT)" "$(CREATE_OUT)" "$(WEB_DIST)" "$(WEB_ASSETS_DIR)"/*

docker:
	docker build --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) --build-arg BUILDTIME="$(BUILDTIME)" -t $(IMG_NAME):$(VERSION) -t $(IMG_NAME):latest $(if $(COMMIT),-t $(IMG_NAME):$(COMMIT)) -f Dockerfile .
