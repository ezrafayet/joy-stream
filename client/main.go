package main

import (
	"fmt"
	"os"

	"keyboard"

	"golang.org/x/term"
)

var title = "     ██╗ ██████╗ ██╗   ██╗███████╗████████╗██████╗ ███████╗ █████╗ ███╗   ███╗\n" + "\r\033[K" +
	"     ██║██╔═══██╗╚██╗ ██╔╝██╔════╝╚══██╔══╝██╔══██╗██╔════╝██╔══██╗████╗ ████║\n" +
	"     ██║██║   ██║ ╚████╔╝ ███████╗   ██║   ██████╔╝█████╗  ███████║██╔████╔██║\n" +
	"██   ██║██║   ██║  ╚██╔╝  ╚════██║   ██║   ██╔══██╗██╔══╝  ██╔══██║██║╚██╔╝██║\n" +
	"╚█████╔╝╚██████╔╝   ██║   ███████║   ██║   ██║  ██║███████╗██║  ██║██║ ╚═╝ ██║\n" +
	" ╚════╝  ╚═════╝    ╚═╝   ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝\n" +
	"Stream a Joycon controller punk ! Even from Japan ! - Client - v0.1.0"

func main() {
	fmt.Println(title)

	source, err := keyboard.NewKeyboard()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer source.Close()

	var keysState = make(map[keyboard.Key]bool)

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
	fmt.Print(returnBeginningLine)
	fmt.Println("Press ESC to exit. Latest event:")
	fmt.Println("")
	fmt.Print(returnBeginningLine)
	for ev := range source.Events() {
		if ev.Key.String() == "ESC" {
			fmt.Print(returnBeginningLine)
			fmt.Println("Exiting... Good bye punk!")
			fmt.Print(returnBeginningLine)
			os.Exit(0)
			break
		}
		keysState[ev.Key] = ev.Type == keyboard.KeyPressed
		fmt.Print(returnBeginningLine)
		fmt.Printf("\r  %s\033[K", keysState)
	}
}
