run: build
	./bin/app
build:
	go build -o bin/app cmd/app/main.go
wire:
	cd ./internal/di && google-wire
docker:
	docker-compose up --build -d