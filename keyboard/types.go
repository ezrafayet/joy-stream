// Package keyboard provides key events (pressed/released) for the current platform.
//
// NewKeyboard() is implemented in exactly one file per build, selected by build tags:
//   - Linux:    evdev_linux.go (evdev, works with Wayland and X11)
//   - Windows / macOS: hook.go (gohook)
// So there is no duplication at runtime — only the right file is compiled for the target OS.
package keyboard

import "fmt"

type KeyEventType int

const (
	KeyPressed KeyEventType = iota
	KeyReleased
)

func (t KeyEventType) String() string {
	switch t {
	case KeyPressed:
		return "pressed"
	case KeyReleased:
		return "released"
	default:
		return "unknown"
	}
}

type KeyEvent struct {
	Key  Key
	Type KeyEventType
}

func (e KeyEvent) String() string {
	return e.Key.String() + " " + e.Type.String()
}

type Key uint16

func (k Key) String() string {
	if name, ok := keyNames[k]; ok {
		return name
	}
	return fmt.Sprintf("key(%d)", k)
}

var keyNames = map[Key]string{
	1:   "ESC",
	// 57:  "SPACE",
	// 103: "UP",
	// 105: "LEFT",
	// 106: "RIGHT",
	// 108: "DOWN",
}

// InputSource emits a stream of key events
type InputSource interface {
	Events() <-chan KeyEvent
	DeviceName() string // e.g. "AT Translated Set 2 keyboard" or "terminal (keyboard)"
	Close() error
}
