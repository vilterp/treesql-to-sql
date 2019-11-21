MAKEFLAGS += -j8

.PHONY: all
all: ui bin/server

.PHONY: ui
ui: ui-deps
	cd ui && yarn build

.PHONY: run
run:
	go run -v main.go

.PHONY: bin/server
bin/server: go-deps
	go build -v -o bin/server .

.PHONY: test
test:
	go test ./...

.PHONY: deps
deps: ui-deps go-deps

.PHONY: ui-deps
ui-deps:
	cd ui && yarn

.PHONY: go-deps
go-deps:
	go mod download
