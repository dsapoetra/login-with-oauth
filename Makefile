.PHONY: migrate-up migrate-down migrate-create

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/migrations -seq $$name

migrate-up:
	migrate -path internal/migrations -database "postgres://:@localhost:5432/poc_oauth?sslmode=disable" up

migrate-down:
	migrate -path internal/migrations -database "postgres://:@localhost:5432/poc_oauth?sslmode=disable" down