include dev.env # .env file

start: build
	@ ./bin/main

build:
	@go build -o ./bin ./cmd/api/main.go

docker-start:
	@docker compose --env-file dev.env up --build

docker-stop:
	@docker-compose rm -v --force --stop
	@docker rmi ticket-booking

docker-init:
	docker run --name ticket-booking -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASSWORD} -p ${DB_PORT}:5432 -d postgres

createdb:
	docker exec -it ticket-booking createdb --username=${DB_USER} ${DB_NAME}

dropdb:
	docker exec -it ticket-booking dropdb ${DB_NAME}