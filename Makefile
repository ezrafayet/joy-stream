build:
	rm -rf build
	go build -o build/server ./server/...

run-server:
	./build/server
