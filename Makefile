build:
	@go build -o build/tasknet cmd/main.go

run: clean build
	@./build/tasknet

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

compose-build:
	@docker compose build

compose-up: compose-build
	@docker compose up -d

compose-up-no-backend: compose-build
	@docker compose up -d --scale tasknet_backend=0 

compose-down:
	@docker compose down

clean:
	@rm -rf ./build/