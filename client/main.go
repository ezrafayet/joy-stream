package main

import (
	"fmt"
	"os"

	"keyboard"

	"golang.org/x/term"
	"github.com/joy-stream/gamepad"
)

var title = 
    "     ██╗ ██████╗ ██╗   ██╗███████╗████████╗██████╗ ███████╗ █████╗ ███╗   ███╗\n" +
	"     ██║██╔═══██╗╚██╗ ██╔╝██╔════╝╚══██╔══╝██╔══██╗██╔════╝██╔══██╗████╗ ████║\n" +
	"     ██║██║   ██║ ╚████╔╝ ███████╗   ██║   ██████╔╝█████╗  ███████║██╔████╔██║\n" +
	"██   ██║██║   ██║  ╚██╔╝  ╚════██║   ██║   ██╔══██╗██╔══╝  ██╔══██║██║╚██╔╝██║\n" +
	"╚█████╔╝╚██████╔╝   ██║   ███████║   ██║   ██║  ██║███████╗██║  ██║██║ ╚═╝ ██║\n" +
	" ╚════╝  ╚═════╝    ╚═╝   ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝\n" +
	"Stream a JoyCon, even from Japan! Client v0.1.2 (versions must align with server)"

func main() {
	fmt.Println(title)

	source, err := keyboard.NewKeyboard()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer source.Close()

	gamepad := gamepad.State{}

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
	fmt.Print(returnBeginningLine)
	fmt.Println("Press ESC to exit. Latest event:")
	fmt.Print(returnBeginningLine)
	fmt.Print("  (no keys yet)\r")
	for ev := range source.Events() {
		if ev.Key.String() == "ESC" {
			fmt.Print(returnBeginningLine)
			fmt.Println("Exiting... Good bye punk!")
			fmt.Print(returnBeginningLine)
			os.Exit(0)
			break
		}
		// keysState[ev.Key] = ev.Type == keyboard.KeyPressed
		fmt.Print(returnBeginningLine)
		fmt.Printf("\r  %s,  Pressed: %s\033[K", gamepad.String(), ev.Key)
	}
}
