dev_api:
	go run cmd/api/main.go

dev_client:
	go run cmd/client/main.go

gen:
	protoc --go_out=. --go-grpc_out=. proto/*

test:
	go test -v ./...
