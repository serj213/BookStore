.PHONY: run, dc, migrate_up, migrate_down, formate, generate

DB_URL = "postgres://postgres:root@localhost:5434/books?sslmode=disable"
PATH_MIGRATION = file://migrations

formate:
	gofmt -w .

migrate_up:
	migrate -source ${PATH_MIGRATION} -database ${DB_URL} up

migrate_down:
	@migrate -source ${PATH_MIGRATION} -database ${DB_URL} down ${VERSION}

dc:
	docker-compose up --remove-orphans --build


run:
	-configPath=config/local.yaml go run cmd/main.go

generate:
	@echo "Generating code from proto files"
	protoc -I docs/protobuf/grpc docs/protobuf/grpc/*.proto --go_out=./pb/grpc --go_opt=paths=source_relative --go-grpc_out=./pb/grpc/ --go-grpc_opt=paths=source_relative