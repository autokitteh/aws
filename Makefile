export GOPRIVATE=github.com/autokitteh/*

GO=go
TAGS=

ARCH=$(shell uname -m)

ifeq ($(COMMIT),)
COMMIT=$(shell git rev-parse HEAD)
endif

ifeq ($(DATE),)
DATE=$(shell date -u "+%Y-%m-%dT%H:%MZ")
endif

ifeq ($(VERSION),)
VERSION=dev
endif

LDFLAGS+=-X 'main.version=${VERSION}' -X 'main.date=${DATE}' -X 'main.commit=${COMMIT}'

ifndef GO_BUILD_OPTS
ifdef DEBUG
GO_BUILD_OPTS+=-gcflags=all="-N -l"
else
GO_BUILD_OPTS=
endif
endif

OUTDIR?=bin
BUILD_OUTDIR=$(OUTDIR)

ifeq (, $(shell which gotestsum))
GOTEST=$(GO) test
else
GOTEST=gotestsum --
endif

define build
$(GO) build --tags "${TAGS}" -o $(BUILD_OUTDIR)/$@ -ldflags="$(LDFLAGS)" $(GO_BUILD_OPTS) ./cmd/$@
endef

define test
$(GOTEST) -v $(GO_TEST_OPTS) -count=1 "$1"
endef

.PHONY: all
all: shellcheck bin lint test

.PHONY: clean
clean:
	rm -fR $(OUTDIR)
	mkdir $(OUTDIR)
	make -C tests clean

.PHONY: bin
bin: awsplugind

.PHONY: build
build:
	$(GO) build $(GO_BUILD_OPTS) ./...
	make lint

.PHONY: debug
debug:
	GO_BUILD_OPTS='-gcflags=all="-N -l"' make bin

.PHONY: awsplugind
awsplugind:
	$(build)

$(OUTDIR)/tools/golangci-lint:
	mkdir -p $(OUTDIR)/tools
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(OUTDIR)/tools" v1.46.2

.PHONY: test
test: test-unit

.PHONY: test-unit
test-unit:
	$(GOTEST) -v --race --tags="unit" -count=1 $(or ${tests},${tests},./...)

.PHONY: lint
lint: $(OUTDIR)/tools/golangci-lint
	$(OUTDIR)/tools/golangci-lint run

.PHONY: shellcheck
shellcheck:
	docker run -v $(shell pwd):/src -w /src koalaman/shellcheck -a -- $(shell find . -name \*.sh)

.PHONY: goreleaser
goreleaser:
	goreleaser release --snapshot --rm-dist

.PHONY: install-githooks
install-githooks:
	./scripts/git-hooks/install.sh
