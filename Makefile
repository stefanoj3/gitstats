VERSION=$(shell git describe || git rev-parse HEAD)
DATE=$(shell date +%s)
LD_FLAGS=-extldflags '-static' -X github.com/stefanoj3/gitstats/internal/cli/cmd.Version=$(VERSION) -X github.com/stefanoj3/gitstats/internal/cli/cmd.BuildTime=$(DATE)

TESTARGS=-v -race -cover -timeout 20s -cpu 24
ifeq ($(CI), true)
TESTARGS=-v -race -coverprofile=coverage.txt -covermode=atomic -timeout 20s -cpu 24
endif

.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: tests
## tests: execute tests
tests:
	@echo "Executing tests"
	@CGO_ENABLED=1 go test $(TESTARGS) ./...

.PHONY: check
## check: run golangci against the codebase
check:
	@echo "Running golangci"
	@docker run -t --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.20.0 golangci-lint run -v

.PHONY: fix
## fix: run goimports against the source code
fix:
	@echo "Running goimports"
	@docker run --rm -v $(PWD):/data cytopia/goimports -w .

.PHONY: fmt
## fmt: run fmt against the source code
fmt:
	@echo "Running fmt"
	@go fmt ./...

.PHONY: build
## build: Builds binary from source
build:
	@echo "Building a new version ${VERSION} - ${DATE}"
	CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "$(LD_FLAGS)" -o dist/gitstats cmd/gitstats/main.go

.PHONY: hook-install
## hook-install: installs git-secrets hook that will run before every commit to help prevent commiting secrest into the repo
# the check is based on simple regex, so it is not bulletproof
hook-install:
	git-secrets --install
	./resources/scripts/git-secret-patterns.sh