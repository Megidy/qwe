start-dependencies:
	@docker-compose up --build -d
run:
	@go run cmd/main.go