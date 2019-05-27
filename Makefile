.PHONY: ui
ui:
	cd ui && yarn build

.PHONY: run
run:
	go run -v main.go

.PHONY: test
test:
	go test ./...
