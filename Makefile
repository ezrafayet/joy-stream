build:
	# go build -o build/client client/main.go
	go build -o build/server server/main.go
	# go build -o build/bridge bridge/joycontrol_bridge.py

run-server:
	./build/server
