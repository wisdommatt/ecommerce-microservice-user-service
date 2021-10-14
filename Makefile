protoc:
	protoc proto/user.proto --go-grpc_out=. --go_out=.
	
run:
	go run server/main.go

tests:
	go test ./... -race -cover

mock:
	go install github.com/golang/mock/mockgen@v1.6.0
	mockgen -source=services/users.go -destination=mocks/users_service.go -package=mocks
	mockgen -source=internal/users/repository.go -destination=mocks/users_repository.go -package=mocks