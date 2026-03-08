package main

import (
	"fmt"
	"os"

	"keyboard"
)

func main() {
	source, err := keyboard.NewKeyboard()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer source.Close()

	fmt.Println("Listening to:", source.DeviceName())
	fmt.Println("Press keys (Ctrl+C to exit). Latest event:")
	for ev := range source.Events() {
		if ev.Key.String() == "ESC" {
			fmt.Println("\nExiting... Good bye punk!")
			os.Exit(0)
		}
		fmt.Printf("\r  %-40s", ev.String())
	}
}
