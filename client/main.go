package main

import (
	"fmt"
	"os"

	"keyboard"
)

func main() {
	source, err := inputsource.NewKeyboard()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer source.Close()

	fmt.Println("Listening to:", source.DeviceName())
	fmt.Println("Press keys (Ctrl+C or close terminal to exit). Events:")
	for ev := range source.Events() {
		fmt.Println("event received:", ev)
		if ev.Key.String() == "ESC" {
			fmt.Println("Exiting...")
			os.Exit(0)
		}
	}
}
