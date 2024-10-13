.PHONY: proto clean clean_omit run test swagger build generate

proto:
	protoc --proto_path=proto \
		--go_out=proto/pb --go_opt=paths=source_relative \
		--go-grpc_out=proto/pb --go-grpc_opt=paths=source_relative \
		proto/*.proto

clean:
	rm -f proto/pb/*.pb.go
	rm -f docs/docs.go docs/swagger.json docs/swagger.yaml

clean_omit:
	find ./proto/pb -name "*.pb.go" -exec sh -c 'sed -i "s/,omitempty//" {}' \;

run:
	go run cmd/main.go

test:
	go test ./... -v

swagger:
	swag init -g cmd/main.go

build:
	go build -o bin/api-gateway cmd/main.go

generate: clean proto swagger clean_omit