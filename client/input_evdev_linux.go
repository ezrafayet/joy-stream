//go:build linux

package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joy-stream/protocol"
)

// Codes Linux input-event-codes.h (EV_KEY)
const (
	evKEY = 1
	keyESC    = 1
	keyA      = 30
	keyS      = 31
	keyD      = 32
	keySpace  = 57
	keyLeft   = 105
	keyRight  = 106
	keyUp     = 103
	keyDown   = 108
)

func runEvdev(conn *net.UDPConn) error {
	dev, err := openKeyboard()
	if err != nil {
		return err
	}
	defer dev.Close()

	var (
		buttons uint16
		dpad    uint8
		quit    bool
		seq     uint16
		mu      sync.Mutex
	)
	pbuf := make([]byte, protocol.PacketSize)
	buf := make([]byte, 24)
	go func() {
		for {
			n, err := dev.Read(buf)
			if err != nil || n != 24 {
				return
			}
			typ := binary.LittleEndian.Uint16(buf[16:18])
			code := binary.LittleEndian.Uint16(buf[18:20])
			val := int32(binary.LittleEndian.Uint32(buf[20:24]))
			if typ != evKEY {
				continue
			}
			mu.Lock()
			down := val == 1 || val == 2 // press or repeat
			switch code {
			case keyESC:
				if val == 0 {
					quit = true
				}
			case keyUp:
				setBit(&dpad, dpadUp, down)
			case keyDown:
				setBit(&dpad, dpadDown, down)
			case keyLeft:
				setBit(&dpad, dpadLeft, down)
			case keyRight:
				setBit(&dpad, dpadRight, down)
			case keySpace:
				setBit16(&buttons, protocol.ButtonA, down)
			case keyA:
				setBit16(&buttons, protocol.ButtonB, down)
			case keyS:
				setBit16(&buttons, protocol.ButtonY, down)
			case keyD:
				setBit16(&buttons, protocol.ButtonX, down)
			}
			// Envoyer tout de suite pour que le relâchement soit visible sans attendre le prochain tick 60 Hz
			p := &protocol.Packet{
				Sequence: seq,
				Buttons:  buttons,
				LX:       128, LY: 128, RX: 128, RY: 128,
				Dpad: dpad,
			}
			seq++
			p.Marshal(pbuf)
			conn.Write(pbuf)
			// Mettre à jour l’affichage tout de suite pour que le relâchement soit visible (pas au prochain tick 60 Hz)
			evdevPrintState(buttons, dpad)
			mu.Unlock()
		}
	}()

	fmt.Println("Contrôles (evdev): Flèches = D-pad, Espace = A, A = B, S = Y, D = X. Échap pour quitter.")
	fmt.Println()

	ticker := time.NewTicker(time.Second / targetHz)
	defer ticker.Stop()

	for range ticker.C {
		mu.Lock()
		if quit {
			mu.Unlock()
			os.Exit(0)
		}
		p := &protocol.Packet{
			Sequence: seq,
			Buttons:  buttons,
			LX:       128, LY: 128, RX: 128, RY: 128,
			Dpad: dpad,
		}
		seq++
		p.Marshal(pbuf)
		conn.Write(pbuf)

		evdevPrintState(buttons, dpad)
		mu.Unlock()
	}
	return nil
}

func evdevPrintState(buttons uint16, dpad uint8) {
	p := &protocol.Packet{Buttons: buttons, Dpad: dpad}
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
}

func setBit(b *uint8, bit uint8, down bool) {
	if down {
		*b |= bit
	} else {
		*b &^= bit
	}
}

func setBit16(b *uint16, bit uint16, down bool) {
	if down {
		*b |= bit
	} else {
		*b &^= bit
	}
}

func openKeyboard() (*os.File, error) {
	for i := 0; i < 32; i++ {
		path := fmt.Sprintf("/dev/input/event%d", i)
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		// Vérifier qu’on a des EV_KEY (lecture non-bloquante impossible sans syscall, on suppose que event0..3 sont souvent clavier/sourice)
		return f, nil
	}
	return nil, fmt.Errorf("aucun périphérique /dev/input/eventX accessible (sudo adduser $USER input ?)")
}
