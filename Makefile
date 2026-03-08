build:
	rm -rf ./build
	go build -o build/server ./server/...
	go build -o build/client ./client/...

build-windows:
	GOOS=windows GOARCH=amd64 go build -o build/client.exe ./client/...

run-server:
	./build/server

run-client:
	sudo ./build/client
