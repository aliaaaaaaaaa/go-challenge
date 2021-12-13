proto:
	protoc api/proto/es.proto   --go_out=plugins=grpc:.

build:
	go mod vendor;
	go build main.go

run:
	go mod vendor;
	go run main.go
