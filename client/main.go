// Client reads keyboard (arrows + Space/A/S/D), sends controller state to the server at 60 Hz.
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/joy-stream/protocol"
)

const (
	defaultPort = 7355
	targetHz   = 60
)

// D-pad bits
const (
	dpadUp    = 1 << 0
	dpadDown  = 1 << 1
	dpadLeft  = 1 << 2
	dpadRight = 1 << 3
)

// Si on ne reçoit plus d’event pour une touche pendant ce temps, on la considère relâchée (la lib n’envoie pas de key release).
const releaseDelay = 600 * time.Millisecond

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("IP du serveur: ")
	serverIP, _ := reader.ReadString('\n')
	serverIP = strings.TrimSpace(serverIP)
	if serverIP == "" {
		fmt.Println("IP requise.")
		os.Exit(1)
	}
	addr := net.UDPAddr{IP: net.ParseIP(serverIP), Port: defaultPort}
	conn, err := net.DialUDP("udp", nil, &addr)
	if err != nil {
		fmt.Printf("Connexion: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	if runtime.GOOS == "linux" {
		if err := runEvdev(conn); err != nil {
			fmt.Println("evdev:", err)
			fmt.Println("Utilisation du clavier terminal (relâchement avec délai).")
			runKeyboard(conn)
		}
		return
	}
	runKeyboard(conn)
}

func runKeyboard(conn *net.UDPConn) {
	if err := keyboard.Open(); err != nil {
		fmt.Printf("Clavier: %v\n", err)
		os.Exit(1)
	}
	defer keyboard.Close()

	fmt.Println("Contrôles: Flèches = D-pad, Espace = A, A = B, S = Y, D = X. Échap pour quitter.")
	fmt.Println()

	var (
		buttons uint16
		dpad    uint8
		mu      sync.Mutex
		last    map[string]time.Time
	)
	last = make(map[string]time.Time)
	keysChan, err := keyboard.GetKeys(10)
	if err != nil {
		fmt.Printf("Clavier GetKeys: %v\n", err)
		os.Exit(1)
	}
	go func() {
		for e := range keysChan {
			if e.Err != nil {
				continue
			}
			now := time.Now()
			mu.Lock()
			switch {
			case e.Key == keyboard.KeyEsc:
				os.Exit(0)
			case e.Key == keyboard.KeyArrowUp:
				dpad |= dpadUp
				last["up"] = now
			case e.Key == keyboard.KeyArrowDown:
				dpad |= dpadDown
				last["down"] = now
			case e.Key == keyboard.KeyArrowLeft:
				dpad |= dpadLeft
				last["left"] = now
			case e.Key == keyboard.KeyArrowRight:
				dpad |= dpadRight
				last["right"] = now
			case e.Rune == ' ':
				buttons |= protocol.ButtonA
				last["A"] = now
			case e.Rune == 'a' || e.Rune == 'A':
				buttons |= protocol.ButtonB
				last["B"] = now
			case e.Rune == 's' || e.Rune == 'S':
				buttons |= protocol.ButtonY
				last["Y"] = now
			case e.Rune == 'd' || e.Rune == 'D':
				buttons |= protocol.ButtonX
				last["X"] = now
			}
			mu.Unlock()
		}
	}()

	seq := uint16(0)
	buf := make([]byte, protocol.PacketSize)
	ticker := time.NewTicker(time.Second / targetHz)
	defer ticker.Stop()

	for range ticker.C {
		mu.Lock()
		now := time.Now()
		for key, t := range last {
			if now.Sub(t) > releaseDelay {
				switch key {
				case "up":
					dpad &^= dpadUp
				case "down":
					dpad &^= dpadDown
				case "left":
					dpad &^= dpadLeft
				case "right":
					dpad &^= dpadRight
				case "A":
					buttons &^= protocol.ButtonA
				case "B":
					buttons &^= protocol.ButtonB
				case "Y":
					buttons &^= protocol.ButtonY
				case "X":
					buttons &^= protocol.ButtonX
				}
				delete(last, key)
			}
		}
		p := &protocol.Packet{
			Sequence: seq,
			Buttons:  buttons,
			LX:       128,
			LY:       128,
			RX:       128,
			RY:       128,
			Dpad:     dpad,
		}
		seq++
		p.Marshal(buf)
		conn.Write(buf)

		if buttons != 0 || dpad != 0 {
			names := p.ButtonNames()
			var dpadStr []string
			if dpad&dpadUp != 0 {
				dpadStr = append(dpadStr, "Haut")
			}
			if dpad&dpadDown != 0 {
				dpadStr = append(dpadStr, "Bas")
			}
			if dpad&dpadLeft != 0 {
				dpadStr = append(dpadStr, "Gauche")
			}
			if dpad&dpadRight != 0 {
				dpadStr = append(dpadStr, "Droite")
			}
			line := "Boutons: " + strings.Join(names, ", ")
			if len(dpadStr) > 0 {
				line += "  |  D-pad: " + strings.Join(dpadStr, ", ")
			}
			fmt.Printf("\r  %s    ", line)
		} else {
			fmt.Print("\r  En attente...    ")
		}
		mu.Unlock()
	}
}
