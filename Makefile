# ==================================================================================== #
# Helpers
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go mod tidy -v
	go run mvdan.cc/gofumpt@latest -w .

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

## lint: Run the golangci linter
.PHONY: lint
lint:
	golangci-lint run

# ==================================================================================== #
## :
## DEVELOPMENT:
# ==================================================================================== #

## test: run the go tests
## : (use `make test pkg=<path-to-package>` to run a specific package, including integrations)
.PHONY: test
test:
	@if [ -z ${pkg} ]; then \
		go test -coverprofile=cover.out -short ./...; \
	else \
		go test -coverprofile=cover.out ${pkg}; \
	fi

## test/coverage: display coverage and indicate if it is less than 80%
.PHONY: test/coverage
test/coverage:
	@go tool cover -func=cover.out | grep "total:" | awk '{print ((int($$3) > 80) != 1) ? "Coverage is less than 80%": "Coverage is greater than 80%"}'

## test/report: display coverage report
test/report:
	@go tool cover -html=cover.out

## serve-docs: generate the godoc documentation and serve it on localhost:6060
.PHONY: serve-docs
serve-docs:
	@godoc -http=:6060 -index -links -v

