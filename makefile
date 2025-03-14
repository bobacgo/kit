run:
	@go run ./examples/main.go

test:
	@go test ./... -v

docs:
	@swag init --parseDependency --parseInternal --dir ./examples --output ./examples/docs