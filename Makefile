# ===============================
# Global Config
# ===============================
BIN_DIR := bin

# ===============================
# Docker Config
# ===============================
DOCKER_REGISTRY   ?= phtom1999
DOCKER_IMAGE_NAME ?= auth-svc
INITIAL_VERSION   ?= v0.0.1
PLATFORM          ?= linux/amd64

# Query the latest semantic version tag from Docker Hub
FETCH_LATEST_TAG = $(shell \
	curl -s -f --connect-timeout 5 "https://hub.docker.com/v2/repositories/$(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME)/tags?page_size=100" 2>/dev/null | \
		grep -o '"name":"[^"]*"' | cut -d'"' -f4 | \
		grep -E '^v?[0-9]+(\.[0-9]+)*$$' | sort -V | tail -n 1)

# Calculate the next version (increments patch version, e.g. v0.0.1 -> v0.0.2)
FETCH_NEXT_TAG = $(shell \
	LATEST="$(FETCH_LATEST_TAG)"; \
	if [ -z "$$LATEST" ]; then \
		echo "$(INITIAL_VERSION)"; \
	else \
		echo "$$LATEST" | awk -F. '{ \
			p = ($$1 ~ /^v/) ? "v" : ""; \
			sub(/^v/, "", $$1); \
			if (NF == 3) printf "%s%s.%s.%d\n", p, $$1, $$2, $$3 + 1; \
			else if (NF == 2) printf "%s%s.%d\n", p, $$1, $$2 + 1; \
			else if (NF == 1) printf "%s%d\n", p, $$1 + 1; \
			else print "$(INITIAL_VERSION)"; \
		}'; \
	fi)

# Only query Docker Hub if a docker-related target is invoked and TAG was not specified
ifeq ($(filter docker-% build-docker push-docker next-tag current-tag,$(MAKECMDGOALS)),)
  TAG ?= $(INITIAL_VERSION)
else
  TAG ?= $(FETCH_NEXT_TAG)
endif

IMAGE        := $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(TAG)
LATEST_IMAGE := $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):latest

# ===============================
# Tool: genrsa
# ===============================
GENRSA_NAME := genrsa
GENRSA_MAIN := ./tools/rsa/main.go
GENRSA_BIN  := $(BIN_DIR)/$(GENRSA_NAME)

## Build the genrsa tool
build-genrsa:
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(GENRSA_BIN) $(GENRSA_MAIN)

## Run the genrsa tool with output directory
run-genrsa:
	@$(GENRSA_BIN) --out-dir=$(out)

# ===============================
# Version Inspection
# ===============================
## View current tag on Docker Hub
current-tag:
	@tag="$(FETCH_LATEST_TAG)"; \
	if [ -z "$$tag" ]; then \
		echo "📦 Current tag on Docker Hub: none (new repository or no semver tags found)"; \
	else \
		echo "📦 Current tag on Docker Hub: $$tag"; \
	fi

## View next tag to be pushed
next-tag:
	@echo "🏷️  Next tag to build/push: $(TAG)"

# ===============================
# Docker Operations
# ===============================
## Build Docker image
docker-build: build-genrsa
	@echo "🔨 Building Docker image: $(IMAGE) (platform: $(PLATFORM))..."
	docker build --platform $(PLATFORM) -t $(IMAGE) -t $(LATEST_IMAGE) .
	@echo "✅ Built: $(IMAGE) and $(LATEST_IMAGE)"

## Push Docker image to registry
docker-push:
	@echo "🚀 Pushing Docker image: $(IMAGE)..."
	docker push $(IMAGE)
	@echo "🚀 Pushing Docker image: $(LATEST_IMAGE)..."
	docker push $(LATEST_IMAGE)
	@echo "🎉 Successfully pushed $(IMAGE) and $(LATEST_IMAGE)!"

## Build and push Docker image
docker-build-push: docker-build docker-push

## Aliases
build-docker: docker-build
push-docker: docker-push

# ===============================
# Clean all
# ===============================
.PHONY: clean build-genrsa run-genrsa current-tag next-tag docker-build docker-push docker-build-push build-docker push-docker

clean:
	rm -rf $(BIN_DIR)
