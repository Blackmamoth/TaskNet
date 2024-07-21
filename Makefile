build:
	@go build -o build/tasknet cmd/main.go

run: clean build
	@./build/tasknet

compose-build:
	@docker compose build

compose-up: compose-build
	@docker compose up -d

compose-down:
	@docker compose down

clean:
	@rm -rf ./build/