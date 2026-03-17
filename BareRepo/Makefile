PROTO_SRC  := proto/repo/repo.proto
PB_OUT     := pkg/pb
SWAGGER_ENTRY := gateway/cmd/swagger.go

.PHONY: generate proto swagger build run-collector run-gateway docker-up

generate: proto swagger

proto:
	protoc \
		--go_out=$(PB_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_OUT) --go-grpc_opt=paths=source_relative \
		-I proto/repo $(PROTO_SRC)

swagger:
	go generate ./gateway/cmd/...

build: generate
	go build -o bin/collector ./collector/cmd
	go build -o bin/gateway   ./gateway/cmd

run-collector:
	go run ./collector/cmd

run-gateway:
	go run ./gateway/cmd

docker-up:
	docker compose up --build
