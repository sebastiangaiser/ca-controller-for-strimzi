# Image configuration
IMAGE_REPO ?= ghcr.io/sebastiangaiser/ca-controller-for-strimzi
IMAGE_TAG ?= dev
IMAGE ?= $(IMAGE_REPO):$(IMAGE_TAG)

# Tool versions
KIND_VERSION ?= v0.27.0
CERT_MANAGER_VERSION ?= v1.17.2

# Kubernetes configuration
KIND_CLUSTER_NAME ?= ca-controller-e2e
NAMESPACE ?= kafka

.PHONY: all
all: build

##@ General

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: build
build: ## Build the Go binary
	go build -o bin/ca-controller-for-strimzi .

.PHONY: test
test: ## Run unit tests
	go test -v ./...

.PHONY: lint
lint: ## Run linter
	golangci-lint run ./...

.PHONY: fmt
fmt: ## Run go fmt
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

##@ Docker

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(IMAGE) .

.PHONY: docker-push
docker-push: ## Push Docker image
	docker push $(IMAGE)

##@ E2E Testing

.PHONY: kind-create
kind-create: ## Create kind cluster for e2e testing
	@if kind get clusters | grep -q "^$(KIND_CLUSTER_NAME)$$"; then \
		echo "Cluster $(KIND_CLUSTER_NAME) already exists"; \
	else \
		kind create cluster --name $(KIND_CLUSTER_NAME) --wait 5m; \
	fi

.PHONY: kind-delete
kind-delete: ## Delete kind cluster
	kind delete cluster --name $(KIND_CLUSTER_NAME)

.PHONY: kind-load
kind-load: docker-build ## Load Docker image into kind cluster
	kind load docker-image $(IMAGE) --name $(KIND_CLUSTER_NAME)

.PHONY: install-cert-manager
install-cert-manager: ## Install cert-manager in the cluster
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/$(CERT_MANAGER_VERSION)/cert-manager.yaml
	kubectl wait --for=condition=Available deployment/cert-manager -n cert-manager --timeout=120s
	kubectl wait --for=condition=Available deployment/cert-manager-webhook -n cert-manager --timeout=120s
	kubectl wait --for=condition=Available deployment/cert-manager-cainjector -n cert-manager --timeout=120s

.PHONY: deploy
deploy: ## Deploy the controller to the cluster
	kubectl create namespace $(NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	helm upgrade --install ca-controller-for-strimzi ./chart \
		--namespace $(NAMESPACE) \
		--set image.repository=$(IMAGE_REPO) \
		--set image.tag=$(IMAGE_TAG) \
		--set image.pullPolicy=Never \
		--wait

.PHONY: undeploy
undeploy: ## Undeploy the controller from the cluster
	helm uninstall ca-controller-for-strimzi --namespace $(NAMESPACE) || true

.PHONY: e2e-setup
e2e-setup: kind-create kind-load install-cert-manager deploy ## Setup e2e test environment

.PHONY: e2e-test
e2e-test: ## Run e2e tests (assumes cluster is ready)
	./hack/e2e-test.sh

.PHONY: e2e
e2e: e2e-setup e2e-test ## Run full e2e test suite

.PHONY: e2e-cleanup
e2e-cleanup: kind-delete ## Cleanup e2e test environment

##@ Release

.PHONY: release
release: ## Create a release (use VERSION=x.y.z)
ifndef VERSION
	$(error VERSION is not set. Use: make release VERSION=x.y.z)
endif
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push origin $(VERSION)
