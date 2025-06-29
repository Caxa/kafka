# Makefile

APP_NAME=go-order-service

.PHONY: all build up down restart logs kafka-send psql init-db create-topic

all: build up

build:
	docker compose build

up:
	docker compose up -d

down:
	docker-compose down -v

restart:
	docker compose down && docker compose up -d --build

logs:
	docker compose logs -f $(APP_NAME)

kafka-send:
	docker exec -it kafka kafka-console-producer \
		--bootstrap-server kafka:9092 \
		--topic orders_topic

psql:
	docker exec -it postgres psql -U order_user -d ordersdb

init-db:
	docker exec -i postgres psql -U order_user -d ordersdb < init.sql

create-topic:
	docker exec -it kafka kafka-topics \
		--create \
		--topic orders_topic \
		--partitions 1 \
		--replication-factor 1 \
		--if-not-exists \
		--bootstrap-server kafka:9092
