include .env/.dev
export
run-dev:
	cd ./go && go run ./cmd/main.go
build-dev:
	cd ./go && go build -o bot.exe ./cmd/main.go
test-dev:
	cd ./go/internal/repository && go test
	cd ./go/internal/service && go test
	cd ./go/internal/bot && go test -run Test_
test-bot:
	cd ./go/internal/bot && go test -run TestRun
compose-up:
	docker compose up -d
compose-down:
	docker compose down
compose-build:
	docker compose build --no-cache
compose-restart:
	docker compose down
	docker compose build --no-cache
	docker compose up -d
gh-action:
	gh workflow run 147279424