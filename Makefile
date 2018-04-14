NAME     := diff2xlsx
VERSION  := v1.2.0
REVISION := $(shell git rev-parse --short HEAD)

SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w \
	-X \"github.com/jiro4989/${NAME}/internal/version.Name=$(NAME)\" \
	-X \"github.com/jiro4989/${NAME}/internal/version.Version=$(VERSION)\" \
	-X \"github.com/jiro4989/${NAME}/internal/version.Revision=$(REVISION)\" \
	-extldflags \"-static\""

MAIN_FILES := cmd/diff2xlsx.go cmd/commands.go

bin/$(NAME): $(SRCS)
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) \
		-o bin/$(NAME) \
		$(MAIN_FILES)

.PHONY: install
install:
	go install $(LDFLAGS) $(MAIN_FILES)

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf vendor/*

.PHONY: cross-build
cross-build: deps
	for os in darwin linux windows; do \
		for arch in amd64 386; do \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$$os-$$arch/$(NAME); \
		done; \
	done

.PHONY: test
test:
	bash ./script/test.sh
