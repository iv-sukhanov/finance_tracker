include .env/.dev
export
run-dev:
	go run cmd/main.go
build-dev:
	go build cmd/main.go
test-dev:
	go test ./internal/repository/*.go
	go test ./internal/service/*.go
test-bot:
	go test ./internal/bot/... -v -run Test_Run