.PHONY: run, sw


run:
	HTTP_ADDR=:8083	go run cmd/main.go

sw:
	swag init -g cmd/main.go --parseDependency --parseInternal