DB_CONTAINER_NAME=postgres_container
DB_NAME=postgres
DB_USER=postgres
DB_PASSWORD=postgres
DB_PORT=5432
DB_IMAGE=postgres:latest
DB_CONN_STR=postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable

export DB_CONN_STR

.PHONY: run db-start db-stop db-clean build

run: build
	DB_CONN_STR=$(DB_CONN_STR) ./app

build:
	go build -o app .

db-start:
	docker run --name $(DB_CONTAINER_NAME) -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_DB=$(DB_NAME) -p $(DB_PORT):5432 -d $(DB_IMAGE)

db-stop:
	docker stop $(DB_CONTAINER_NAME) || true
	docker rm $(DB_CONTAINER_NAME) || true

db-clean: db-stop
	docker rmi $(DB_IMAGE) || true

