include .env
export PROJECT_ROOT=$(shell pwd)

pg-up:
	@docker compose up -d postgres

pg-down:
	@docker compose down postgres

pg-clean-up:
	@read -p "Очистить pg volume? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down postgres && \
	  	rm -rf out/pgdata && \
	  	echo "Успешно очищено"; \
	else \
  	  	echo "Очистка отменена"; \
  	fi

migrate-pg-create:
	@if [ -z "$(seq)" ]; then \
  		echo "Отсутствует seq"; \
  		exit 1; \
	fi; \
	docker compose run --rm postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@if [ -z "$(action)" ]; then \
  		echo "Отсутствует action"; \
  		exit 1; \
	fi; \
	docker compose run --rm postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable \
		$(action)

migrate-pg-up:
	@make migrate-action action=up

migrate-pg-down:
	@make migrate-action action=down

port-forward:
	@docker compose up -d port-forwarder

port-forward-close:
	@docker compose down port-forwarder

