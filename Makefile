PROGRAM_NAME = gophermart

.PHONY: help dep fmt test

dep: ## Get the dependencies
	go mod download

fmt: ## Format the source files
	gofumpt -l -w .

test: dep ## Run tests
	go test -timeout 5m -race -covermode=atomic -coverprofile=.coverage.out ./... && \
	go tool cover -func=.coverage.out | tail -n1 | awk '{print "Total test coverage: " $$3}'
	@rm .coverage.out

cover: dep ## Run app tests with coverage report
	go test -timeout 5m -race -covermode=atomic -coverprofile=.coverage.out ./... && \
	go tool cover -html=.coverage.out -o .coverage.html
	## Open coverage report in default system browser
	xdg-open .coverage.html
	## Remove coverage report
	sleep 2 && rm -f .coverage.out .coverage.html

build:
	go build -o ./cmd/gophermart/gophermart ./cmd/gophermart

lint: lint/sources lint/openapi ## Run all linters

lint/sources: ## Lint the source files
	golangci-lint run --timeout 5m
	govulncheck ./...

lint/openapi: ## Lint openapi specifications
	@echo "Lint OpenAPI specifications"
	@for spec in $(OPENAPI_SPECS) ; do echo "* lint $$spec"; vacuum lint -t -q -x $$spec ; done