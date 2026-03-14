package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"keyboard"
	"network"

	"github.com/joy-stream/gamepad"
	"golang.org/x/term"
)

var title = 
    "     ██╗ ██████╗ ██╗   ██╗███████╗████████╗██████╗ ███████╗ █████╗ ███╗   ███╗\n" +
	"     ██║██╔═══██╗╚██╗ ██╔╝██╔════╝╚══██╔══╝██╔══██╗██╔════╝██╔══██╗████╗ ████║\n" +
	"     ██║██║   ██║ ╚████╔╝ ███████╗   ██║   ██████╔╝█████╗  ███████║██╔████╔██║\n" +
	"██   ██║██║   ██║  ╚██╔╝  ╚════██║   ██║   ██╔══██╗██╔══╝  ██╔══██║██║╚██╔╝██║\n" +
	"╚█████╔╝╚██████╔╝   ██║   ███████║   ██║   ██║  ██║███████╗██║  ██║██║ ╚═╝ ██║\n" +
	" ╚════╝  ╚═════╝    ╚═╝   ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝\n" +
	"Stream a JoyCon, even from Japan! Client v0.1.2 (versions must align with server)"

// Mapping used for keyboard
type Mapping struct {
	DpadUp       uint16
	DpadDown     uint16
	DpadLeft     uint16
	DpadRight    uint16
	StickUp      uint16
	StickDown    uint16
	StickLeft    uint16
	StickRight   uint16
	TriggerLeft  uint16
	TriggerRight uint16
}

func main() {
	fmt.Println(title)

	//////////////////// server address /////////////////
	fmt.Print("Server address (host or host:port) [like 127.0.0.1:7355]: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "no input")
		os.Exit(1)
	}
	serverAddr := strings.TrimSpace(scanner.Text())
	if serverAddr == "" {
		panic("no server address provided")
	} else if !strings.Contains(serverAddr, ":") {
		panic("server address must be in the format host:port")
	}
	fmt.Println("Using server:", serverAddr)
	////////////////////////////////////////////////////

	sender, err := network.NewSender(serverAddr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "network:", err)
		os.Exit(1)
	}
	defer sender.Close()

	source, err := keyboard.NewKeyboard()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer source.Close()

	var gamepadMu sync.Mutex
	var seq uint16
	gp := gamepad.State{}
	mapping := Mapping{
		DpadUp:       23, // i
		DpadDown:     37, // k
		DpadLeft:     36, // j
		DpadRight:    38, // l
		StickUp:      17, // w
		StickDown:    31, // s
		StickLeft:    30, // a
		StickRight:   32, // d
		TriggerLeft:  16, // q
		TriggerRight: 18, // e
	}

	// Now switch to raw mode so keypresses aren't echoed
	if term.IsTerminal(int(os.Stdin.Fd())) {
		state, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Fprintln(os.Stderr, "terminal:", err)
			os.Exit(1)
		}
		defer term.Restore(int(os.Stdin.Fd()), state)
	}

	returnBeginningLine := "\r\033[K"

	fmt.Println("Listening to:", source.DeviceName())
	if os.Getenv("XDG_SESSION_TYPE") == "wayland" {
		fmt.Println("Detected Wayland — gohook needs X11. Run: GDK_BACKEND=x11 ./build/client")
	} else {
		fmt.Println("(Linux: if no keys are detected, run under X11: GDK_BACKEND=x11 or use an X11 session)")
	}
	// 60 Hz send loop
	go func() {
		fmt.Fprintf(os.Stderr, "[client] sending to %s at 60 Hz\n", serverAddr)
		ticker := time.NewTicker(time.Second / 60)
		defer ticker.Stop()
		for range ticker.C {
			gamepadMu.Lock()
			gp.Sequence = seq
			seq++
			buf := make([]byte, gamepad.PacketSize)
			gp.Marshal(buf)
			gamepadMu.Unlock()
			if err := sender.Send(buf); err != nil {
				fmt.Fprintf(os.Stderr, "send error: %v\n", err)
			}
		}
	}()

	fmt.Print(returnBeginningLine)
	fmt.Println("Press ESC to exit. Latest event:")
	fmt.Print(returnBeginningLine)
	gamepadMu.Lock()
	fmt.Printf("\r%s\033[K", gp.String())
	gamepadMu.Unlock()
	for ev := range source.Events() {
		if ev.Key.String() == "ESC" {
			fmt.Print(returnBeginningLine)
			fmt.Println("Exiting... Good bye punk!")
			fmt.Print(returnBeginningLine)
			os.Exit(0)
			break
		}
		gamepadMu.Lock()
		if ev.Type == keyboard.KeyPressed {
			switch uint16(ev.Key) {
			case mapping.DpadUp:
				gp.SetDpadUp(true)
			case mapping.DpadDown:
				gp.SetDpadDown(true)
			case mapping.DpadLeft:
				gp.SetDpadLeft(true)
			case mapping.DpadRight:
				gp.SetDpadRight(true)
			case mapping.StickUp:
				gp.SetStickUp(true)
			case mapping.StickDown:
				gp.SetStickDown(true)
			case mapping.StickLeft:
				gp.SetStickLeft(true)
			case mapping.StickRight:
				gp.SetStickRight(true)
			case mapping.TriggerLeft:
				gp.SetTriggerLeft(true)
			case mapping.TriggerRight:
				gp.SetTriggerRight(true)
			}
		} else if ev.Type == keyboard.KeyReleased {
			switch uint16(ev.Key) {
			case mapping.DpadUp:
				gp.SetDpadUp(false)
			case mapping.DpadDown:
				gp.SetDpadDown(false)
			case mapping.DpadLeft:
				gp.SetDpadLeft(false)
			case mapping.DpadRight:
				gp.SetDpadRight(false)
			case mapping.StickUp:
				gp.SetStickUp(false)
			case mapping.StickDown:
				gp.SetStickDown(false)
			case mapping.StickLeft:
				gp.SetStickLeft(false)
			case mapping.StickRight:
				gp.SetStickRight(false)
			case mapping.TriggerLeft:
				gp.SetTriggerLeft(false)
			case mapping.TriggerRight:
				gp.SetTriggerRight(false)
			}
		}
		stateStr := gp.String()
		gamepadMu.Unlock()
		fmt.Print(returnBeginningLine)
		fmt.Printf("\r%s,  Pressed: %s\033[K", stateStr, ev.Key)
	}
}
