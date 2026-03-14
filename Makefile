.PHONY: build build-windows run-server run-client

build:
	rm -rf ./build
	go build -o build/server ./server
	go build -o build/client ./client

# Cross-compile client for Windows (from Linux). Requires: sudo apt-get install gcc-mingw-w64
build-windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
	CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ \
	go build -o build/client.exe ./client/...

run-server:
	./build/server

run-client:
	sudo ./build/client

test:
	go test ./gamepad/...
