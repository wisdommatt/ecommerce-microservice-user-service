protoc:
	protoc proto/user.proto --go-grpc_out=. --go_out=.
	
run:
	go run server/main.go

test:
	go test ./... -race -cover