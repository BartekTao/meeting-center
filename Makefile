GO_CMD=go

MAIN_FILE=./cmd/server.go

all:
	$(GO_CMD) run $(MAIN_FILE)

db-up:
	docker-compose -f ./deployment/Docker-compose.yml up -d
db-down:
	docker-compose -f ./deployment/Docker-compose.yml down
	
build:
	$(GO_CMD) build -o meeting-room-api $(MAIN_FILE)

run:
	./meeting-room-api

gen:
	go get github.com/99designs/gqlgen@v0.17.45
	go run github.com/99designs/gqlgen generate