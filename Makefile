protoc:
	protoc proto/user.proto --
	
run:
	go run server/main.go