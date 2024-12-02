include common.mk

PACKAGES=$(shell go list ./...)
BUILDDIR?=$(CURDIR)/build
OUTPUT?=$(BUILDDIR)/intellixd

MODULE_NAME := github.com/IntelliXLabs/intellix
HTTPS_GIT := https://$(MODULE_NAME).git
CGO_ENABLED ?= 1

# Process Docker environment varible TARGETPLATFORM
# in order to build binary with correspondent ARCH
# by default will always build for linux/amd64
TARGETPLATFORM ?=
GOOS ?= linux
GOARCH ?= amd64
GOARM ?=

ifeq (linux/arm,$(findstring linux/arm,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=arm
	GOARM=7
endif

ifeq (linux/arm/v6,$(findstring linux/arm/v6,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=arm
	GOARM=6
endif

ifeq (linux/arm64,$(findstring linux/arm64,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=arm64
	GOARM=7
endif

ifeq (linux/386,$(findstring linux/386,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=386
endif

ifeq (linux/amd64,$(findstring linux/amd64,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=amd64
endif

ifeq (linux/mips,$(findstring linux/mips,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=mips
endif

ifeq (linux/mipsle,$(findstring linux/mipsle,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=mipsle
endif

ifeq (linux/mips64,$(findstring linux/mips64,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=mips64
endif

ifeq (linux/mips64le,$(findstring linux/mips64le,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=mips64le
endif

ifeq (linux/riscv64,$(findstring linux/riscv64,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=riscv64
endif

#? all: Run target build, test and install
all: build test install
.PHONY: all

include tests.mk

###############################################################################
###                                Build Intellixd                           ###
###############################################################################

#? build: Build Intellixd
build:
	go build -mod=readonly -ldflags "-s -w" -o $(OUTPUT) ./cmd/intellixd/
.PHONY: build

build-debug:
	go build -gcflags="all=-N -l" --trimpath=false -o $(OUTPUT) ./cmd/intellixd/
.PHONY: build

#? install: Install Intellixd to GOBIN
install:
	CGO_ENABLED=$(CGO_ENABLED) go install $(BUILD_FLAGS) -tags $(BUILD_TAGS) ./cmd/intellixd
.PHONY: install

install-lib:
	bash scripts/install-lib.sh

install-wasm-testdata:
	bash scripts/install-wasm-testdata.sh



###############################################################################
###                              Distribution                               ###
###############################################################################

#? go-mod-cache: Download go modules to local cache
go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	go mod download
.PHONY: go-mod-cache

#? go.sum: Ensure dependencies have not been modified
go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	go mod verify
	go mod tidy

#? draw_deps: Generate deps graph
draw_deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	goviz -i $(MODULE_NAME)/cmd/intellixd -d 3 | dot -Tpng -o dependency-graph.png
.PHONY: draw_deps

get_deps_bin_size:
	@# Copy of build recipe with additional flags to perform binary size analysis
	$(eval $(shell go build -work -a $(BUILD_FLAGS) -tags $(BUILD_TAGS) -o $(OUTPUT) ./cmd/intellixd/ 2>&1))
	@find $(WORK) -type f -name "*.a" | xargs -I{} du -hxs "{}" | sort -rh | sed -e s:${WORK}/::g > deps_bin_size.log
	@echo "Results can be found here: $(CURDIR)/deps_bin_size.log"
.PHONY: get_deps_bin_size


###############################################################################
###                  Formatting, linting, and vetting                       ###
###############################################################################

format:
	find . -name '*.go' -type f -not -path "*.git*" -not -name '*.pb.go' -not -name '*pb_test.go' | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "*.git*"  -not -name '*.pb.go' -not -name '*pb_test.go' | xargs goimports -w -local $(MODULE_NAME)
.PHONY: format

#? lint: Run latest golangci-lint linter
lint:
	@echo "--> Running linter"
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3 run
.PHONY: lint

#? vulncheck: Run latest govulncheck
vulncheck:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
.PHONY: vulncheck

#? lint-typo: Run codespell to check typos
lint-typo:
	which codespell || pip3 install codespell
	@codespell
.PHONY: lint-typo

#? lint-typo: Run codespell to auto fix typos
lint-fix-typo:
	@codespell -w
.PHONY: lint-fix-typo

DESTINATION = ./index.html.md


###############################################################################
###                           Documentation                                 ###
###############################################################################

#? check-docs-toc: Verify that important design docs have ToC entries.
check-docs-toc:
	@./docs/presubmit.sh
.PHONY: check-docs-toc

###############################################################################
###                            Docker image                                 ###
###############################################################################

# On Linux, you may need to run `DOCKER_BUILDKIT=1 make build-docker` for this
# to work.
#? build-docker: Build docker image intellix/intellix
build-docker:
	docker build \
		--label=intellix \
		--tag="intellix/intellix" \
		-f DOCKER/Dockerfile .
.PHONY: build-docker

###############################################################################
###                       Local testnet using docker                        ###
###############################################################################

#? build-linux: Build linux binary on other platforms
build-linux:
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) $(MAKE) build
.PHONY: build-linux

#? build-docker-localnode: Build the "localnode" docker image
build-docker-localnode:
	@cd networks/local && make
.PHONY: build-docker-localnode

# Implements test splitting and running. This is pulled directly from
# the github action workflows for better local reproducibility.

GO_TEST_FILES != find $(CURDIR) -name "*_test.go"

# default to four splits by default
NUM_SPLIT ?= 2

$(BUILDDIR):
	mkdir -p $@

# The format statement filters out all packages that don't have tests.
# Note we need to check for both in-package tests (.TestGoFiles) and
# out-of-package tests (.XTestGoFiles).
$(BUILDDIR)/packages.txt:$(GO_TEST_FILES) $(BUILDDIR)
	go list -f "{{ if (or .TestGoFiles .XTestGoFiles) }}{{ .ImportPath }}{{ end }}" ./... | grep -v "^$(MODULE_NAME)/test/" | grep -v "^$(MODULE_NAME)/internal/test/" | grep -v "^$(MODULE_NAME)/_back/" | sort > $@

split-test-packages:$(BUILDDIR)/packages.txt
	split -d -n l/$(NUM_SPLIT) $< $<.
test-group-%:split-test-packages
	cat $(BUILDDIR)/packages.txt.$* | xargs go test -mod=readonly -timeout=15m -race -coverprofile=$(BUILDDIR)/$*.profile.out

test-in-ci:$(BUILDDIR)/packages.txt install-wasm-testdata install-lib
	cat $(BUILDDIR)/packages.txt | xargs go test -mod=readonly -timeout=15m -race -coverprofile=$(BUILDDIR)/coverage.txt

test-runtime:
	@if [ ! -f "lib/runtime.h" ]; then \
		echo "Error: lib/runtime.h not found"; \
		echo "Please build runtime from github.com/IntelliXLabs/iwasm and copy it to lib/"; \
		exit 1; \
	fi
	@if ! ls lib/libruntime.* >/dev/null 2>&1; then \
		echo "Error: libruntime shared library not found in lib/ directory"; \
		echo "Expected one of: libruntime.so, libruntime.dylib, or libruntime.dll"; \
		echo "Please build runtime from github.com/IntelliXLabs/iwasm and copy it to lib/"; \
		exit 1; \
	fi
	LD_LIBRARY_PATH=$(PWD)/lib CGO_LDFLAGS=-L$(PWD)/lib go test ./tests/...
.PHONY: test_runtime


#? help: Get more info on make commands.
help: Makefile
	@echo " Choose a command run in comebft:"
	@sed -n 's/^#?//p' $< | column -t -s ':' |  sort | sed -e 's/^/ /'
.PHONY: help

include mk-docker.mk
