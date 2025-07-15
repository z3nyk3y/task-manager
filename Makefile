-include .env


up: ## up task-manager and database in deamon.
	docker compose up -d

up-no-detached: ## up task-manager and database and do not detache from std out.
	docker compose up

build-app: ## build app
	docker compose build app

migrations-apply: ## applies migration. Use MGR_NUM_UP in .env file to configurate how many migrations need to apply. Defaullt all migrations will be set.
	docker run --rm --name task-manager-migrator -v ./migrations:/migrations --network task-manager_task-manager_net migrate/migrate:4 \
	-path=/migrations/ -database "pgx5://$(DB_LOGIN):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up $(MGR_NUM_UP)

migrations-rollback: ## applies migration. Use MGR_NUM_UP in .env file to configurate how many migrations need to rollback. Defaullt all migrations will be rollback.
	docker run --rm --name task-manager-migrator -v ./migrations:/migrations --network task-manager_task-manager_net migrate/migrate:4 \
	-path=/migrations/ -database "postgres://$(DB_LOGIN):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down $(MGR_NUM_DOWN)

help:
	@echo "available commands:"; \
	grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
