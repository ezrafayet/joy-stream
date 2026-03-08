//go:build linux

package inputsource

import (
	"fmt"
	"os"

	"github.com/holoplot/go-evdev"
)

// evdevSource implements InputSource by reading from a Linux evdev device.
type evdevSource struct {
	dev  *evdev.InputDevice
	name string
	ch   chan KeyEvent
}

// NewKeyboard opens the first available keyboard (evdev device with EV_KEY
// and letter keys like KEY_A). Skips touchpads/mice that also have EV_KEY.
func NewKeyboard() (InputSource, error) {
	paths, err := evdev.ListDevicePaths()
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		if _, err := os.Open("/dev/input/event0"); err != nil && os.IsPermission(err) {
			return nil, fmt.Errorf("permission denied on /dev/input/* — run: sudo adduser %s input, then log out and log back in", os.Getenv("USER"))
		}
		return nil, fmt.Errorf("no input devices found under /dev/input")
	}
	for _, p := range paths {
		dev, err := evdev.OpenWithFlags(p.Path, os.O_RDONLY)
		if err != nil {
			if os.IsPermission(err) {
				return nil, fmt.Errorf("permission denied on %s — run: sudo adduser %s input, then log out and log back in", p.Path, os.Getenv("USER"))
			}
			continue
		}
		types := dev.CapableTypes()
		hasKey := false
		for _, t := range types {
			if t == evdev.EV_KEY {
				hasKey = true
				break
			}
		}
		if !hasKey {
			_ = dev.Close()
			continue
		}
		// Prefer a real keyboard: must have KEY_A (touchpads have BTN_* but not letter keys)
		codes := dev.CapableEvents(evdev.EV_KEY)
		hasLetter := false
		for _, c := range codes {
			if c == evdev.KEY_A {
				hasLetter = true
				break
			}
		}
		if !hasLetter {
			_ = dev.Close()
			continue
		}
		name, _ := dev.Name()
		path := dev.Path()
		if name == "" {
			name = path
		} else {
			name = name + " (" + path + ")"
		}
		ch := make(chan KeyEvent, 64)
		s := &evdevSource{dev: dev, name: name, ch: ch}
		go s.run()
		return s, nil
	}
	return nil, fmt.Errorf("no keyboard device found (try: sudo adduser %s input, then log out and back in)", os.Getenv("USER"))
}

func (s *evdevSource) Events() <-chan KeyEvent {
	return s.ch
}

func (s *evdevSource) DeviceName() string {
	return s.name
}

func (s *evdevSource) Close() error {
	err := s.dev.Close()
	close(s.ch)
	return err
}

func (s *evdevSource) run() {
	for {
		e, err := s.dev.ReadOne()
		if err != nil {
			return
		}
		if e.Type != evdev.EV_KEY {
			continue
		}
		ev := KeyEvent{Key: Key(e.Code)}
		switch e.Value {
		case 0:
			ev.Type = KeyReleased
		case 1, 2:
			ev.Type = KeyPressed
		default:
			continue
		}
		select {
		case s.ch <- ev:
		default:
			// channel full, drop
		}
	}
}
