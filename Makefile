.PHONY: all
all: ui bin/server

.PHONY: ui
ui:
	cd ui && yarn build

.PHONY: run
run:
	go run -v main.go

.PHONY: bin/server
bin/server:
	go build -v -o bin/server .

.PHONY: test
test:
	go test ./...

.PHONY: deps
deps:
	cd ui && yarn
	go mod download
