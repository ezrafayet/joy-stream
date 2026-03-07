build:
	rm -rf ./build
	go build -o build/server ./server/...
	go build -o build/client ./client/...

run-server:
	./build/server

run-client:
	./build/client
