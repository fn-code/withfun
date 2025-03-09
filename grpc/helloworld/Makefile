clean:
	rm pb/*

gen:
	protoc --proto_path=proto proto/*.proto  --go_out=:pb --go-grpc_out=:pb --grpc-gateway_out=:pb

gen2:
	protoc --go_out=:pb --go_opt=paths=source_relative \
        --go-grpc_out=:pb --go-grpc_opt=paths=source_relative \
        --proto_path=proto proto/*.proto

.PHONY: clean gen gen2