# go utils

# variables
# Binary Name
GO_BIN_NAME == ???

# Path to operate on
GO_FSPATH == ???

# Path to build single binary to.
GO_BUILD_OUT_FSPATH = ???

# Packages to operate one
GO_PKG_LIST = ???

# Path to build all binaries to
GO_BUILD_OUT_ALL_FSPATH = ???

# Supported Platforms
GO_PLATFORMS=darwin windows

# Supported Architectures (on other platforms than linux)
GO_ARCHITECTURES=386 amd64

# Other supported platforms
GO_LINUX_PLATFORMS=linux

# Supported Linux Architectures
GO_LINUX_ARCHITECTURES=386 amd64 arm arm64

## Print
go-print: 
	@echo
	@echo -- GO --
	@echo GO_FSPATH: 				$(GO_FSPATH)
	@echo GO_PKG_LIST: 				$(GO_PKG_LIST)
	@echo GO_BUILD_OUT_FSPATH: 		$(GO_BUILD_OUT_FSPATH)
	@echo GO_PLATFORMS:             $(GO_PLATFORMS) linux
	@echo GO_ARCHITECTURES:         $(GO_ARCHITECTURES)
	@echo GO_LINUX_ARCHITECTURES    $(GO_LINUX_ARCHITECTURES)
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
	cd $(GO_FSPATH) && go build -v -o $(GO_BUILD_OUT_FSPATH)/$(GO_BIN_NAME) .

## Remove build
go-build-clean:
	@echo "Removing build"
	rm -rf $(GO_BUILD_OUT_FSPATH)/$(GO_BIN_NAME)

## Cross-compile to Supported Archs and OSes
go-build-all:
	@echo "Building All Supported Architectures & Platforms"
	mkdir -p $(GO_BUILD_OUT_ALL_FSPATH)
	cd $(GO_FSPATH)
	$(foreach GOOS, $(GO_PLATFORMS),\
	$(foreach GOARCH, $(GO_ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); \
		go build -v -o $(GO_BUILD_OUT_ALL_FSPATH)/$(GO_BIN_NAME)-$(GOOS)-$(GOARCH) $(GO_FSPATH))))
	$(foreach GOOS, $(GO_LINUX_PLATFORMS),\
	$(foreach GOARCH, $(GO_LINUX_ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); \
		go build -v -o $(GO_BUILD_OUT_ALL_FSPATH)/$(GO_BIN_NAME)-$(GOOS)-$(GOARCH) $(GO_FSPATH))))

## Clean / remove cross-compile directory
go-build-clean-all:
	@echo "Cleaning all builds"
	rm -rf $(GO_BUILD_OUT_ALL_FSPATH)

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




