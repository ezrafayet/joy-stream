//go:build darwin || windows

package inputsource

import (
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
)

// keyboardSource implements InputSource using the terminal keyboard package (darwin/windows).
type keyboardSource struct {
	ch   chan KeyEvent
	name string
}

// NewKeyboard opens the terminal keyboard and returns an InputSource.
// On Darwin/Windows we use the terminal; key release may not be reported.
func NewKeyboard() (InputSource, error) {
	if err := keyboard.Open(); err != nil {
		return nil, err
	}
	ch := make(chan KeyEvent, 64)
	name := "terminal keyboard"
	if os.Getenv("TERM") != "" {
		name = fmt.Sprintf("terminal keyboard (%s)", os.Getenv("TERM"))
	}
	s := &keyboardSource{ch: ch, name: name}
	go s.run()
	return s, nil
}

func (s *keyboardSource) Events() <-chan KeyEvent {
	return s.ch
}

func (s *keyboardSource) DeviceName() string {
	return s.name
}

func (s *keyboardSource) Close() error {
	keyboard.Close()
	close(s.ch)
	return nil
}

func (s *keyboardSource) run() {
	keysChan, err := keyboard.GetKeys(10)
	if err != nil {
		return
	}
	for e := range keysChan {
		if e.Err != nil {
			continue
		}
		key := keyboardToKey(e)
		if key == 0 {
			continue
		}
		ev := KeyEvent{Key: key, Type: KeyPressed}
		select {
		case s.ch <- ev:
		default:
		}
	}
}

// keyboardToKey maps keyboard package Key/Rune to our Key (evdev-like codes for consistent String()).
func keyboardToKey(e keyboard.KeyEvent) Key {
	if e.Key != 0 {
		switch e.Key {
		case keyboard.KeyEsc:
			return 1
		case keyboard.KeyArrowUp:
			return 103
		case keyboard.KeyArrowDown:
			return 108
		case keyboard.KeyArrowLeft:
			return 105
		case keyboard.KeyArrowRight:
			return 106
		case keyboard.KeySpace:
			return 57
		default:
			return Key(e.Key)
		}
	}
	switch e.Rune {
	case ' ':
		return 57
	case 'a', 'A':
		return 30
	case 's', 'S':
		return 31
	case 'd', 'D':
		return 32
	default:
		return Key(e.Rune)
	}
}
