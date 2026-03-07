// Package main runs the UDP input server: listens for controller packets,
// tracks clients, and displays server IP and connected clients' joycon state.
package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/joy-stream/protocol"
)

const (
	listenAddr     = ":7355"       // UDP port for controller packets
	displayRefresh = 500 * time.Millisecond
	clientTimeout  = 3 * time.Second // consider client gone after no packet
)

type clientState struct {
	Addr     string
	Last     *protocol.Packet
	LastSeq  uint16
	LastSeen time.Time
}

func main() {
	conn, err := net.ListenPacket("udp", listenAddr)
	if err != nil {
		fmt.Printf("Failed to start UDP server: %v\n", err)
		return
	}
	defer conn.Close()

	// Show server IP(s) so clients can connect
	printServerIPs(listenAddr)

	clients := make(map[string]*clientState)
	var mu sync.RWMutex

	// Goroutine: read UDP packets and update client state
	go func() {
		buf := make([]byte, 1024)
		for {
			n, remote, err := conn.ReadFrom(buf)
			if err != nil {
				return
			}
			if n != protocol.PacketSize {
				continue
			}
			p, err := protocol.ParsePacket(buf[:n])
			if err != nil {
				continue
			}

			addr := remote.String()
			mu.Lock()
			c, ok := clients[addr]
			if !ok {
				c = &clientState{Addr: addr}
				clients[addr] = c
			}
			// Ignore old sequence numbers (out-of-order or duplicate)
			if p.Sequence >= c.LastSeq || c.Last == nil {
				c.Last = p
				c.LastSeq = p.Sequence
			}
			c.LastSeen = time.Now()
			mu.Unlock()
		}
	}()

	// Goroutine: periodically remove stale clients and refresh display
	ticker := time.NewTicker(displayRefresh)
	defer ticker.Stop()
	for range ticker.C {
		mu.Lock()
		now := time.Now()
		for addr, c := range clients {
			if now.Sub(c.LastSeen) > clientTimeout {
				delete(clients, addr)
			}
		}
		printClients(clients)
		mu.Unlock()
	}
}

func printServerIPs(port string) {
	_, portNum, _ := net.SplitHostPort(port)
	if portNum == "" {
		portNum = "7355"
	}

	fmt.Println("--- Joy-Stream UDP Input Server ---")
	fmt.Println("Server is listening for controller packets (UDP).")
	fmt.Println()

	// Public IP = celle que ton frère au Japon (ou n'importe où) peut utiliser
	publicIP := fetchPublicIP()
	if publicIP != "" {
		fmt.Printf("  >>> Connect from anywhere (e.g. Japan): %s:%s\n", publicIP, portNum)
		fmt.Println()
	}

	// IPs locales pour le même réseau
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		shown := make(map[string]struct{})
		for _, a := range addrs {
			ipNet, ok := a.(*net.IPNet)
			if !ok || ipNet.IP.IsLoopback() || ipNet.IP.To4() == nil {
				continue
			}
			ip := ipNet.IP.String()
			if _, ok := shown[ip]; ok {
				continue
			}
			shown[ip] = struct{}{}
			// Skip Docker bridges for cleaner output
			if strings.HasPrefix(ip, "172.17.") || strings.HasPrefix(ip, "172.18.") {
				continue
			}
			fmt.Printf("  On local network: %s:%s\n", ip, portNum)
		}
	}

	fmt.Println("------------------------------------")
}

func fetchPublicIP() string {
	client := &http.Client{Timeout: 3 * time.Second}
	// Services qui renvoient juste l'IP en texte (pas de JSON)
	urls := []string{
		"https://api.ipify.org",
		"https://ifconfig.me/ip",
	}
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

func printClients(clients map[string]*clientState) {
	var line string
	if len(clients) == 0 {
		line = fmt.Sprintf("[%s] No clients connected. Waiting for packets...", time.Now().Format("15:04:05"))
	} else {
		line = fmt.Sprintf("[%s] %d client(s): ", time.Now().Format("15:04:05"), len(clients))
		for _, c := range clients {
			line += c.Addr
			if c.Last != nil {
				btns := c.Last.ButtonNames()
				dpad := c.Last.DpadNames()
				if len(btns) > 0 || len(dpad) > 0 {
					line += " → pressé: "
					if len(btns) > 0 {
						line += "boutons " + strings.Join(btns, ",")
					}
					if len(dpad) > 0 {
						if len(btns) > 0 {
							line += "  "
						}
						line += "D-pad " + strings.Join(dpad, ",")
					}
				}
			}
			line += "  "
		}
	}
	fmt.Print("\r\033[K" + line)
}
