//go:build windows || darwin

// Package keyboard uses github.com/robotn/gohook for global key events.
// Used on Windows and macOS only. On Linux, evdev is used instead (see evdev_linux.go).
// macOS: enable Accessibility permission for the app.
package keyboard

import (
	"fmt"
	"os"
	"sync"

	"github.com/robotn/gohook"
)

type hookSource struct {
	ch   chan KeyEvent
	once sync.Once
}

// hookSource implements InputSource
var _ InputSource = (*hookSource)(nil)

// NewKeyboard starts a global keyboard hook and returns an InputSource that
// emits KeyPressed and KeyReleased. Uses gohook (X11 on Linux, not Wayland).
func NewKeyboard() (InputSource, error) {
	ch := make(chan KeyEvent, 64)
	s := &hookSource{ch: ch}
	go s.run()
	return s, nil
}

func (s *hookSource) Events() <-chan KeyEvent {
	return s.ch
}

func (s *hookSource) DeviceName() string {
	return "keyboard (gohook)"
}

func (s *hookSource) Close() error {
	s.once.Do(func() {
		hook.End()
	})
	return nil
}

func (s *hookSource) run() {
	evChan := hook.Start()
	defer close(s.ch)
	first := true
	for ev := range evChan {
		switch ev.Kind {
		case hook.KeyDown, hook.KeyHold:
			s.ch <- KeyEvent{Key: Key(ev.Keycode), Type: KeyPressed}
		case hook.KeyUp:
			s.ch <- KeyEvent{Key: Key(ev.Keycode), Type: KeyReleased}
		}
		if first && (ev.Kind == hook.KeyDown || ev.Kind == hook.KeyUp || ev.Kind == hook.KeyHold) {
			first = false
			fmt.Fprintln(os.Stderr, "keyboard: hook receiving events (if you see no keys in the app, check key codes)")
		}
	}
}
