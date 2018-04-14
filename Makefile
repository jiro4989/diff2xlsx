NAME     := diff2xlsx
VERSION  := v1.3.4
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
cross-build:
	rm -rf $(DIST_DIR)
	bash ./script/cross-build.sh $(NAME) $(VERSION) $(LDFLAGS) $(MAIN_FILES)

.PHONY: test
test:
	bash ./script/test.sh

.PHONY: release
release:
	-rm $(DIST_DIR)/*.tar.gz
	ls -d $(DIST_DIR)/* | while read -r d; do cp $(COPY_FILES) $$d/; done
	ls -d $(DIST_DIR)/* | while read -r d; do tar czf $$d.tar.gz $$d; done
	ghr $(VERSION) $(DIST_DIR)/
	go install cmd/*
