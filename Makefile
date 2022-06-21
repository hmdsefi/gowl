GO               = go
M                = $(shell printf "\033[34;1m>>\033[0m")

# Check richgo does exist.
ifeq (, $(shell which richgo))
$(warning "could not find richgo in $(PATH), run: go get github.com/kyoh86/richgo")
endif

.PHONY: test sync codecov test-app

.PHONY: default
default: all

.PHONY: all
all: test

.PHONY: test
test: sync
	$(info running tests)
	 richgo test -v ./...

.PHONY: codecov
codecov: sync
	$(info running tests coverage)
	sh build/script/coverage.sh

.PHONY: sync
sync: 
	$(info downloading dependencies)
	go get -v ./...

.PHONY: fmt
fmt:
	$(info $(M) format code)
	@ret=0 && for d in $$($(GO) list -f '{{.Dir}}' ./... | grep -v /vendor/); do \
		$(GO) fmt $$d/*.go || ret=$$? ; \
		done ; exit $$ret

.PHONY: lint
lint: ## Run linters
	$(info $(M) running golangci linter)
	golangci-lint run --timeout 5m0s ./...

