
# This make file uses composition to keep things KISS and easy.
# In the boilerpalte make files dont do any includes, because you will create multi permutations of possibilities.



# git include


include ./boilerplate/core/help.mk

include ./boilerplate/core/bs.mk
include ./boilerplate/core/os.mk
include ./boilerplate/core/gitr.mk
include ./boilerplate/core/go.mk


# examples of how to override in root make file
override GO_FSPATH = $(PWD)
override GO_BUILD_OUT_FSPATH = $(GOPATH)/bin/bs
override BS_ROOT_FSPATH = XXX
GO_ARCH=go-arch
override GO_ARCH=go-arch_override


STATIK_DEST = $(PWD)/statiks

.PHONY: help this-statiks this-scan-statiks this-build

## Print all settings
this-print:
	$(MAKE) bs-print

	$(MAKE) os-print
	
	$(MAKE) gitr-print

	$(MAKE) go-print

## Example to Print Variable override from make
this-print-ex:
	# prints specific overides
	@echo BS_ROOT_FSPATH: 	$(BS_ROOT_FSPATH)

	@echo GO_ARCH: 	$(GO_ARCH)
	
## Example to Print Variable override from env
this-print-env-ex:
	# prints specific overides
	# Example call of override from env variable
	# ``` GO_ARCH=GO_ARCH_FROMENV make -e print-env ```
	@echo BS_ROOT_FSPATH: 	$(BS_ROOT_FSPATH)

	@echo GO_ARCH: 	$(GO_ARCH)


## Build this.
this-build: this-dep this-statiks this-scan-statiks
	$(MAKE) go-build

## Delete the build.
this-build-clean: this-dep-clean
	rm -rf $(GOPATH)/bin/bs

	# delete all generated stuff
	rm -rf $(PWD)/statiks

this-dep:
	# add binaries needs to build
	go get -u github.com/rakyll/statik

this-dep-clean:
	rm -rf $(GOPATH)/bin/statik

this-statiks:
	@statik -src=$(PWD)/boilerplate -ns bp -p bp -dest=$(STATIK_DEST) -f

