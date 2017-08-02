NAME     := invoker
VERSION  := v0.0.1
REVISION := $(shell git rev-parse --short HEAD)

SRCS    := $(shell find . -type f -name '*.go')
OPTIONS := -tags netgo -installsuffix netgo
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

DIST_DIRS := find * -type d -exec

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	go build $(OPTIONS) $(LDFLAGS) -o bin/$(NAME)

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf vendor/*

.PHONY: deps
deps: glide
	# dep ensure
	glide install

.PHONY: dep
dep:
	go get -u github.com/golang/dep/cmd/dep

.PHONY: glide
glide:
ifeq ($(shell command -v glide 2> /dev/null),)
	curl https://glide.sh/get | sh
endif

.PHONY: install
install:
	go install $(LDFLAGS)

.PHONY: test
test:
	go test -cover -v `go list ./... | grep -v /vendor/`

.PHONY: ci-test
ci-test:
	echo "" > coverage.txt
	for d in $$(go list ./... | grep -v /vendor/); do \
		go test -coverprofile=profile.out -covermode=atomic -race -v $$d; \
		if [ -f profile.out ]; then \
			cat profile.out >> coverage.txt; \
			rm profile.out; \
		fi; \
	done

.PHONY: update-deps
update-deps: glide
	# dep ensure -update
	glide install
