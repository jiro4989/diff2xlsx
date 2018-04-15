NAME     := diff2xlsx
VERSION  := v1.3.9
REVISION := $(shell git rev-parse --short HEAD)

SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w \
	-X \"github.com/jiro4989/$(NAME)/internal/version.Name=$(NAME)\" \
	-X \"github.com/jiro4989/$(NAME)/internal/version.Version=$(VERSION)\" \
	-X \"github.com/jiro4989/$(NAME)/internal/version.Revision=$(REVISION)\" \
	-extldflags \"-static\""

# mainパッケージ帰属のソース
MAIN_FILES := cmd/diff2xlsx.go cmd/commands.go

# 配布物に含めるファイル
COPY_FILES := README.md CHANGELOG.md

# 配布物の出力先
DIST_DIR   := dist/$(VERSION)

build: $(SRCS)
	go build -a -tags netgo -installsuffix netgo $(LDFLAGS) \
		-o bin/$(NAME) \
		$(MAIN_FILES)
	go install cmd/*

test: build
	bash ./script/test.sh

install:
	go install $(LDFLAGS) $(MAIN_FILES)

clean:
	rm -rf bin/*
	rm -rf vendor/*

deps:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/tcnksm/ghr
	dep ensure

cross-build:
	-rm -rf $(DIST_DIR)
	bash ./script/cross-build.sh $(NAME) $(VERSION) $(LDFLAGS) $(MAIN_FILES)

archive: cross-build
	ls -d $(DIST_DIR)/* | while read -r d; do cp $(COPY_FILES) $$d/; done
	bash ./script/arch.sh $(DIST_DIR)

release: archive
	ghr $(VERSION) $(DIST_DIR)/

.PHONY: build test install clean deps cross-build archive release
