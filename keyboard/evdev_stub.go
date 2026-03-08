//go:build !linux && !darwin && !windows

package keyboard

import "fmt"

// NewKeyboard returns an error on non-Linux (evdev is Linux-only).
func NewKeyboard() (InputSource, error) {
	return nil, fmt.Errorf("keyboard input source only supported on Linux")
}
