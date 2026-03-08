package keyboard

import "fmt"

// KeyEventType is press or release.
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

// KeyEvent is a single key press or release.
type KeyEvent struct {
	Key  Key
	Type KeyEventType
}

func (e KeyEvent) String() string {
	return e.Key.String() + " " + e.Type.String()
}

// Key is a key code (e.g. from evdev). Use String() for display.
type Key uint16

// String returns a readable name for common keys, otherwise "key(N)".
func (k Key) String() string {
	if name, ok := keyNames[k]; ok {
		return name
	}
	return fmt.Sprintf("key(%d)", k)
}

// keyNames maps evdev key codes to short names (subset).
var keyNames = map[Key]string{
	1:   "ESC",
	30:  "A",
	31:  "S",
	32:  "D",
	57:  "SPACE",
	103: "UP",
	105: "LEFT",
	106: "RIGHT",
	108: "DOWN",
}

// InputSource emits a stream of key events. Easy to fake in tests; works for
// keyboard today, controller tomorrow; avoids callback spaghetti.
type InputSource interface {
	Events() <-chan KeyEvent
	DeviceName() string // e.g. "AT Translated Set 2 keyboard" or "terminal (keyboard)"
	Close() error
}
