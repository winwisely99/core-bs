# go utils

# variables
# Path to operate on
GO_FSPATH == ???

# Path to build binary to.
GO_BUILD_OUT_FSPATH = ???

# Packages to operate one
GO_PKG_LIST = ???

## Print
go-print: 
	@echo
	@echo -- GO --
	@echo GO_FSPATH: 				$(GO_FSPATH)
	@echo GO_PKG_LIST: 				$(GO_PKG_LIST)
	@echo GO_BUILD_OUT_FSPATH: 		$(GO_BUILD_OUT_FSPATH)
	@echo


## Boilerplate is updated from the Boostrap repo ( REDUNDANT )
go-boilerplate-update:
	# See: https://github.com/lyft/boilerplate
	# Example: See: https://github.com/lyft/flytepropeller/tree/master/boilerplate
	# TODO: This will be redundant once we have BS releases working.
	@boilerplate/update.sh

## Build the code
go-build:
	@echo Building
	cd $(GO_FSPATH) && go build -v -o $(GO_BUILD_OUT_FSPATH) .

## Run the code
go-run:
	@echo Running
	cd $(GO_FSPATH) && go run -v .

## Format with go-fmt
go-fmt:
	@echo Formatting
	cd $(GO_FSPATH) && go fmt .

## Lint with golangci-lint
go-lint:
	@echo Linting
	cd $(GO_FSPATH) && golangci-lint run --no-config --issues-exit-code=0 --timeout=5m

## Run the tests
go-test:
	@echo Running tests
	cd $(GO_FSPATH) && go test -race -v .

## Run the tests with coverage
go-test-coverage:
	@echo Running tests with coverage
	cd $(GO_FSPATH) && go test -short -coverprofile cover.out -covermode=atomic ${GO_PKG_LIST}

## Display test coverage
go-display-coverage:
	@echo Displaying test coverage
	cd $(GO_FSPATH) && go tool cover -html=cover.out




