package gamepad

import (
	"encoding/binary"
	"errors"
)

// PacketSize is the wire size in bytes. Layout is defined here (joycon), not by any external protocol.
const PacketSize = 12

// Wire layout (big-endian):
//  0-1: Sequence (uint16)
//  2:   Dpad     (uint8: bit0=Up, bit1=Down, bit2=Left, bit3=Right)
//  3:   Buttons  (uint8: bit6=ZL, bit7=ZR; other bits reserved)
//  4:   LX       (uint8: 0=left, 128=center, 255=right)
//  5:   LY       (uint8: 0=up, 128=center, 255=down)
//  6-7: RX, RY   (uint8: right stick, 128=center)
//  8-11: Reserved
const (
	wireSeq    = 0
	wireDpad   = 2
	wireButtons = 3
	wireLX     = 4
	wireLY     = 5
	wireRX     = 6
	wireRY     = 7
)

// Dpad bits in wire byte 2.
const (
	DpadUp    = 1 << 0
	DpadDown  = 1 << 1
	DpadLeft  = 1 << 2
	DpadRight = 1 << 3
)

// Button bits in wire byte 3 (for triggers; extend as needed).
const (
	ButtonZL = 1 << 6
	ButtonZR = 1 << 7
)

// State is the joycon controller state. This struct defines the logical model; the wire format above encodes it.
type State struct {
	Sequence     uint16
	DpadUp       bool
	DpadDown     bool
	DpadLeft     bool
	DpadRight    bool
	StickUp      bool
	StickDown    bool
	StickLeft    bool
	StickRight   bool
	TriggerLeft  bool
	TriggerRight bool
}

// Marshal encodes s into dst (at least PacketSize bytes). Returns the slice written, or nil if dst is too short.
func (s *State) Marshal(dst []byte) []byte {
	if len(dst) < PacketSize {
		return nil
	}
	binary.BigEndian.PutUint16(dst[wireSeq:wireSeq+2], s.Sequence)
	var dpad, btn uint8
	if s.DpadUp {
		dpad |= DpadUp
	}
	if s.DpadDown {
		dpad |= DpadDown
	}
	if s.DpadLeft {
		dpad |= DpadLeft
	}
	if s.DpadRight {
		dpad |= DpadRight
	}
	if s.TriggerLeft {
		btn |= ButtonZL
	}
	if s.TriggerRight {
		btn |= ButtonZR
	}
	dst[wireDpad] = dpad
	dst[wireButtons] = btn
	dst[wireLX] = stickByte(s.StickLeft, s.StickRight)
	dst[wireLY] = stickByte(s.StickUp, s.StickDown)
	dst[wireRX] = 128
	dst[wireRY] = 128
	dst[8], dst[9], dst[10], dst[11] = 0, 0, 0, 0
	return dst[:PacketSize]
}

func stickByte(neg, pos bool) uint8 {
	if neg {
		return 0
	}
	if pos {
		return 255
	}
	return 128
}

// Unmarshal decodes a PacketSize-byte packet into a new State. The wire structure is defined in this package.
func Unmarshal(data []byte) (*State, error) {
	if len(data) != PacketSize {
		return nil, errors.New("packet must be 12 bytes")
	}
	s := &State{}
	s.Sequence = binary.BigEndian.Uint16(data[wireSeq : wireSeq+2])
	dpad := data[wireDpad]
	btn := data[wireButtons]
	s.DpadUp = dpad&DpadUp != 0
	s.DpadDown = dpad&DpadDown != 0
	s.DpadLeft = dpad&DpadLeft != 0
	s.DpadRight = dpad&DpadRight != 0
	s.TriggerLeft = btn&ButtonZL != 0
	s.TriggerRight = btn&ButtonZR != 0
	lx := data[wireLX]
	ly := data[wireLY]
	s.StickLeft = lx < 128
	s.StickRight = lx > 128
	s.StickUp = ly < 128
	s.StickDown = ly > 128
	return s, nil
}

func (s *State) SetDpadUp(pressed bool)   { s.DpadUp = pressed }
func (s *State) SetDpadDown(pressed bool)  { s.DpadDown = pressed }
func (s *State) SetDpadLeft(pressed bool) { s.DpadLeft = pressed }
func (s *State) SetDpadRight(pressed bool) { s.DpadRight = pressed }
func (s *State) SetStickUp(pressed bool)   { s.StickUp = pressed }
func (s *State) SetStickDown(pressed bool) { s.StickDown = pressed }
func (s *State) SetStickLeft(pressed bool) { s.StickLeft = pressed }
func (s *State) SetStickRight(pressed bool) { s.StickRight = pressed }
func (s *State) SetTriggerLeft(pressed bool) { s.TriggerLeft = pressed }
func (s *State) SetTriggerRight(pressed bool) { s.TriggerRight = pressed }

