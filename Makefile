.PHONY: ui
ui:
	cd ui && yarn build

.PHONY: run
run:
	go run main.go

.PHONY: test
test:
	go test ./...
