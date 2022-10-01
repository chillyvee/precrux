VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT  := $(shell git log -1 --format='%H')

all: install

LD_FLAGS = -X github.com/chillyvee/precrux/cmd/precrux/cmd.Version=$(VERSION) \
	-X github.com/chillyvee/precrux/cmd/precrux/cmd.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

build:
	@go build -mod readonly $(BUILD_FLAGS) -o build/ ./cmd/precrux/...

install:
	@go install -mod readonly $(BUILD_FLAGS) ./cmd/precrux/...

build-linux:
	@GOOS=linux GOARCH=amd64 go build --mod readonly $(BUILD_FLAGS) -o ./build/precrux ./cmd/precrux

test:
	@go test -timeout 20m -mod readonly -v ./...

clean:
	rm -rf build

build-precrux-docker:
	docker build -t chillyvee/precrux:$(VERSION) -f ./docker/precrux/Dockerfile .

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(dir $(mkfile_path))

.PHONY: all lint test race msan tools clean build
