include .env
export

export PROJECT_ROOT = $(CURDIR)


env-up:
	@docker-compose up -d dokkee-postgres

env-down:
	@docker-compose down dokkee-postgres

env-ps:
	@docker-compose ps

env-clean: 
	docker-compose down dokkee-postgres
	@cmd /c "if exist out\pgdata (rmdir /s /q out\pgdata && echo Cleaned up out/pgdata)\
	 else echo out/pgdata does not exist, nothing to clean."

env-port-forward:
	@docker-compose up -d port-forwarder

env-port-close:
	@docker-compose down port-forwarder

migrate-create:
ifeq ($(seq),)
	@echo Need to write seq. Usage: make migrate-create seq=name
	@exit 1
else
	docker-compose run --rm dokkee-postgres-migrate \
	create \
	-ext sql \
	-dir /migrations \
	-seq "$(seq)"
endif

migrate-up:
	make migrate-action action=up

migrate-down:
	make migrate-action action=down

migrate-action:
ifeq ($(action),)
	@echo Need to write action. Usage: make migrate-action action=up/down
	@exit 1
else
	@docker-compose run --rm dokkee-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@dokkee-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"
endif
