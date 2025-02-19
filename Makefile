run: build
	./bin/app
build:
	go build -o bin/app cmd/app/main.go
wire:
	cd ./internal/di && wire
docker-up:
	docker compose up --build
swag:
	swag init --exclude docker,nginx,assets,pkg --md ./docs --parseInternal --parseDependency --parseDepth 2 -g cmd/app/main.go
test:
	go test ./internal/api/handler/...
mock:
	go generate ./internal/domain/interfaces