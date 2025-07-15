.PHONY: all test

help:  ## Show this help.
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-30s\033[0m %s\n", $$1, $$2}'

build:  ## Build the project.
	@go build .

test:  ## Run tests. Needs bash > 4.0, gotestsum, panicparse, and script
	@which gotestsum > /dev/null || go install gotest.tools/gotestsum@latest
	@which panicparse > /dev/null || go install github.com/maruel/panicparse/v2@latest
	@mkdir -p coverage
	@/opt/homebrew/bin/bash -c "CGO_ENABLED=0 GOTRACEBACK=all script -q /dev/null doppler --project server --config github run -- gotestsum --hide-summary=skipped --format-hide-empty-pkg -- -coverprofile=coverage/cover.out -short -vet=all -tags test -shuffle=on ./... |& panicparse -rel-path"

staticcheck:  ## Run staticcheck
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@staticcheck ./...

mcp-inspector:  ## Run mcp-inspector
	$(MAKE) build
	@npx @modelcontextprotocol/inspector ./chalk-mcp

.DEFAULT:
	@echo Unknown command. Available commands below
	@echo
	@make help
