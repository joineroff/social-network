include .env
export $(shell sed 's/=.*//' .env)

COMPOSE = docker-compose -f docker-compose.${ENV}.yml
LOGS= docker-compose -f docker-compose.${ENV}.yml logs -f --tail=100

.PHONY:mysql-migrate-up
mysql-migrate-up:
	@docker run --rm \
		-v $(shell pwd)/mysql/migrations:/migrations \
		--network host migrate/migrate \
		-path=/migrations/ \
		-database mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@/${MYSQL_DATABASE} \
		up

.PHONY:mysql-migrate-force
mysql-migrate-force:
	@docker run --rm \
		-v $(shell pwd)/mysql/migrations:/migrations \
		--network host migrate/migrate \
		-path=/migrations/ \
		-database mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@/${MYSQL_DATABASE} \
		force 1

.PHONY:mysql-migrate-down
mysql-migrate-down:
	@docker run --rm \
		-v $(shell pwd)/mysql/migrations:/migrations \
		--network host migrate/migrate \
		-path=/migrations/ \
		-database mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@/${MYSQL_DATABASE} \
		down 1
generate_backend_cfg:
	@ ./bin/generate-backend-config.sh

generate_haproxy_cfg:
	@ ./bin/generate-haproxy-config.sh

up:
	@$(COMPOSE) up -d

start: generate_backend_cfg up

stop:
	@$(COMPOSE) stop

down:
	@$(COMPOSE) down

rebuild-backend:
	@$(COMPOSE) build --no-cache backend

restart: stop start

logs-%:
	@${LOGS} $*

status:
	@echo "-----------------------\n"
	@$(COMPOSE) ps
	@echo "-----------------------\n"

.PHONY:backend-run
backend-run:
	cd backend; \
	go run ./main.go

.PHONY:backend-run-race
backend-run-race:
	go run --race  ./cmd/backend

.PHONY:backend-lint
backend-lint:
	@docker run --rm \
	-v $(PWD):/app \
	-w /app \
	golangci/golangci-lint \
	golangci-lint run -v --timeout 5m

.PHONY:backend-test
backend-test:
	go test -v ./...

.PHONY:backend-coverage
backend-coverage:
	go test -cover ./...

.PHONY:backend-gosec
backend-gosec:
	@docker run --rm -v $(PWD):/app -w /app securego/gosec /app/...

.DEFAULT_GOAL=backend-run
