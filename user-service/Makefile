.PHONY: proto clean run test migrations-up migrations-down build generate

proto:
	protoc --proto_path=proto \
		--go_out=proto/pb --go_opt=paths=source_relative \
		--go-grpc_out=proto/pb --go-grpc_opt=paths=source_relative \
		proto/*.proto

clean:
	rm -f proto/pb/*.pb.go

run:
	go run cmd/main.go

test:
	go test ./... -v

migrations-up:
	psql -d $(DB_NAME) -U $(DB_USER) -f migrations/user_service_up.sql

migrations-down:
	psql -d $(DB_NAME) -U $(DB_USER) -f migrations/user_service_down.sql

build:
	go build -o bin/user-service cmd/main.go

generate: clean proto
