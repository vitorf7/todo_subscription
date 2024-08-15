GOPATH:=$(shell go env GOPATH)

# ----- Installing -----

.PHONY: install
install:
	go mod download

# ----- Linting -----

.PHONY: lint
lint:
	golangci-lint run

.PHONY: ci-lint
ci-lint:
	@ [ -e ./bin/golangci-lint ] || wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest
	./bin/golangci-lint run --timeout=10m

# ----- Testing -----

BUILDENV := CGO_ENABLED=0
TESTFLAGS := -short -cover --coverprofile=coverage.out

.PHONY: test
test:
	go generate ./...
	GOPRIVATE=github.com/utilitywarehouse/* $(BUILDENV) go test $(TESTFLAGS) ./...
# ----- Building the binary -----

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
base_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))
SERVICE ?= $(base_dir)

.PHONY: clean
clean:
	rm -f $(SERVICE)

GIT_SHA := $(GITHUB_SHA)
ifeq ($(GIT_SHA),)
  GIT_SHA := $(shell git rev-parse HEAD)
endif

LINKFLAGS := -ldflags '-s -w ${xflags} -extldflags "-static"'
.PHONY: build
build:
	cd $(CURDIR); CGO_ENABLED=0 go build $(LINKFLAGS) -o $(SERVICE)

$(SERVICE): build

.PHONY: all
all: clean $(LINTER) lint test build


golangci-lint:
	golangci-lint run ./... --timeout=10m

FILES_TO_FIX_COUNT=$(shell golines ${FILE_PATHS} --list-files --base-formatter=gofumpt --ignore-generated --max-len=140 | wc -l)
golines-lint:
	if [ $(FILES_TO_FIX_COUNT) != 0 ]; then golines ${FILE_PATHS} --list-files --base-formatter=gofumpt --ignore-generated --max-len=140 && exit 1; fi


gofumpt-format:
	gofumpt -l -w .

FILE_PATHS=$(shell git diff --name-only --diff-filter=d -- "*.go")
ifeq ($(FILE_PATHS),)
	FILE_PATHS=./
endif
golines-fix:
	golines ${FILE_PATHS} -w --base-formatter="gofumpt" --ignore-generated --max-len=140

.PHONY: format
format: gofumpt-format golines-fix ## format Go Code

.PHONY: gqlgen-generate
gqlgen-generate:
	go run github.com/99designs/gqlgen generate

# ----- Running in DEV -----

RUN_PARAMS :=

.PHONY: dev
dev:
	go run $(LINKFLAGS) $(CURDIR)  $(RUN_PARAMS)
