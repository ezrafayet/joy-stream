// Package main runs the UDP input server: listens for controller packets
// and prints gamepad state when received.
package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/joy-stream/gamepad"
	"udp"
)

const listenAddr = ":7355"

func main() {
	recv, err := udp.NewReceiver(listenAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start UDP server: %v\n", err)
		return
	}
	defer recv.Close()

	fmt.Printf("Joy-Stream UDP server listening on %s\n", listenAddr)

	var mu sync.Mutex
	clients := make(map[string]*gamepad.State)

	recv.Start(func(data []byte, from net.Addr) {
		if len(data) != gamepad.PacketSize {
			return
		}
		st, err := gamepad.Unmarshal(data)
		if err != nil {
			return
		}
		addr := from.String()
		mu.Lock()
		clients[addr] = st
		mu.Unlock()
		fmt.Print("\r\033[K")
		fmt.Printf("[%s] %s: %s", time.Now().Format("15:04:05"), addr, st.String())
	})

	select {}
}
