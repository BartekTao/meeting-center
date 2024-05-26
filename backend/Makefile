GO_CMD=go

MAIN_FILE=./cmd/api/server.go

all:
	docker-compose -f ./deployment/Docker-compose.yml up -d

down:
	docker-compose -f ./deployment/Docker-compose.yml down

db-up:
	docker-compose -f ./deployment/Docker-compose-mongodb.yml up -d
db-down:
	docker-compose -f ./deployment/Docker-compose-mongodb.yml down
	
build:
	$(GO_CMD) build -o meeting-room-api $(MAIN_FILE)

run:
	./meeting-room-api

gen:
	go get github.com/99designs/gqlgen@v0.17.45
	go run github.com/99designs/gqlgen generate
