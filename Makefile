include .env/.dev
export
run-dev:
	go run cmd/main.go
build-dev:
	go build cmd/main.go