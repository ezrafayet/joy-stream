// Package main runs the UDP input server: listens for controller packets
// and prints gamepad state when received.
package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joy-stream/gamepad"
	"udp"
)

const listenAddr = ":7355"
const port = "7355"

func main() {
	recv, err := udp.NewReceiver(listenAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start UDP server: %v\n", err)
		return
	}
	defer recv.Close()

	fmt.Printf("Joy-Stream UDP server listening on %s\n", listenAddr)
	if publicIP := fetchPublicIP(); publicIP != "" {
		fmt.Printf("Connect from anywhere: %s:%s\n", publicIP, port)
		fmt.Println("(Forward UDP 7355 on your router to this machine; allow it in the firewall.)")
	}
	fmt.Println()

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

func fetchPublicIP() string {
	client := &http.Client{Timeout: 3 * time.Second}
	urls := []string{"https://api.ipify.org", "https://ifconfig.me/ip"}
	for _, u := range urls {
		resp, err := client.Get(u)
		if err != nil {
			continue
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}
		b, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}
		ip := strings.TrimSpace(string(b))
		if ip != "" && net.ParseIP(ip) != nil {
			return ip
		}
	}
	return ""
}
