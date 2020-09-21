GO_ENV ?= CGO_ENABLED=0 

THIS_FILE := $(lastword $(MAKEFILE_LIST))

tools:
	@cd backend \
		&& go get \
			google.golang.org/grpc \
			github.com/golang/protobuf/protoc-gen-go \
			github.com/bufbuild/buf/cmd/buf \
			github.com/google/wire/cmd/wire
	@go get -d github.com/envoyproxy/protoc-gen-validate
	@cd $(go env GOPATH)/src/github.com/envoyproxy/protoc-gen-validate \
		&& make build

regen:
	@find . -type f -name '*.pb.*.go' -o -name '*.pb.go' -delete
	@for PROTO in $(shell find . -type f -name '*.proto' | grep -v proto/google/api) ; do \
		echo "Compiling" $${PROTO} ; \
		protoc \
			-I=. \
			-I=$(shell go env GOPATH)/src/ \
			-I=$(shell go env GOPATH)/src/github.com/envoyproxy/protoc-gen-validate \
			--go_out=plugins=grpc,paths=source_relative:./ \
			--validate_out=lang=go,paths=source_relative:. \
			$${PROTO}; \
	done;
	@cd backend \
		&& buf check lint || :;
	@wire gen ./...

bin/%:
	@cd backend \
		&& echo "building cmd/$(shell basename $@)" \
		&& $(GO_ENV) go build \
			-trimpath \
			-gcflags='-e -l' \
			-ldflags='-w -s -extldflags "-static" -X main.version=${VERSION} -X main.gitCommit=$(GIT_COMMIT)' \
			-o $@ \
			./cmd/$(shell basename $@)
bin:
	@for CMD in $(shell find ./backend/cmd/* -maxdepth 0 -type d -name '*' ) ; do \
		$(MAKE) -f $(THIS_FILE) bin/$$(basename $${CMD}); \
	done;