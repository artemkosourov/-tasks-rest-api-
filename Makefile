.PHONY: run build clean docker-up docker-down docker-build docker-logs docker-ps restart clean-all status app-logs db-logs up-with-migrate quick-start

run:
	go run cmd/app/main.go

build:
	go build -o bin/tasks-api cmd/app/main.go

clean:
	rm -rf bin/

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-ps:
	docker-compose ps

app-logs:
	docker-compose logs -f app

db-logs:
	docker-compose logs -f postgres

restart: docker-down docker-up

clean-all: docker-down
	docker-compose rm -f
	docker volume rm tasks_postgres_data 2>/dev/null || true
	docker network rm tasks_network 2>/dev/null || true

status:
	@echo "Статус контейнеров:"
	@docker-compose ps
	@echo ""
	@echo "Использование портов:"
	@-lsof -i :8080 -i :5433 2>/dev/null || echo "Порты свободны"

up-with-migrate: docker-up
	@echo "Контейнеры запущены"
	@echo "Миграции выполняются автоматически при старте app..."
	@echo "Проверьте логи: make app-logs"
	@echo "API: http://localhost:8080"

quick-start: docker-build docker-up
	@echo "Server is running"
	@echo "API: http://localhost:8080"
	@echo "Logs: make app-logs"
