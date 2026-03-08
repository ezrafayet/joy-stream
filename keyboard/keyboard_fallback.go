//go:build darwin || windows

package keyboard

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
)

// keyboardSource implements InputSource using the terminal keyboard package (darwin/windows).
// The terminal API only reports key press (and repeat when held); we emit a synthetic
// release per key 50ms after the last press event for that key, so long press doesn't release early.
type keyboardSource struct {
	ch           chan KeyEvent
	name         string
	releaseDelay time.Duration
	pending      map[Key]time.Time // one release per key, at this time
	mu           sync.Mutex
	done         chan struct{}
}

// NewKeyboard opens the terminal keyboard and returns an InputSource.
// On Darwin/Windows the terminal only gives key press (and repeat when held); we emit
// KeyReleased 50ms after the last press event for each key.
func NewKeyboard() (InputSource, error) {
	if err := keyboard.Open(); err != nil {
		return nil, err
	}
	ch := make(chan KeyEvent, 64)
	name := "terminal keyboard"
	if os.Getenv("TERM") != "" {
		name = fmt.Sprintf("terminal keyboard (%s)", os.Getenv("TERM"))
	}
	s := &keyboardSource{
		ch:           ch,
		name:         name,
		releaseDelay: 50 * time.Millisecond,
		pending:      make(map[Key]time.Time),
		done:         make(chan struct{}),
	}
	go s.run()
	go s.releaseLoop()
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
	close(s.done)
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
		select {
		case s.ch <- KeyEvent{Key: key, Type: KeyPressed}:
		default:
		}
		s.scheduleRelease(key)
	}
}

// scheduleRelease sets release for this key at now+releaseDelay (replaces any previous for same key).
func (s *keyboardSource) scheduleRelease(key Key) {
	s.mu.Lock()
	s.pending[key] = time.Now().Add(s.releaseDelay)
	s.mu.Unlock()
}

func (s *keyboardSource) releaseLoop() {
	tick := time.NewTicker(20 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-s.done:
			return
		case <-tick.C:
		}
		s.mu.Lock()
		now := time.Now()
		for key, at := range s.pending {
			if !now.Before(at) {
				delete(s.pending, key)
				select {
				case s.ch <- KeyEvent{Key: key, Type: KeyReleased}:
				default:
				}
			}
		}
		s.mu.Unlock()
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
