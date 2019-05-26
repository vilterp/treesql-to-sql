.PHONY: ui
ui:
	cd ui && yarn build

.PHONY: run
run: ui
	go run -v main.go

.PHONY: test
test:
	go test ./...
