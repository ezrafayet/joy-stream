package gamepad

import (
	"encoding/binary"
	"errors"

	"fmt"
)

const PacketSize = 5

// Wire layout (big-endian):
//  0-1: Sequence (uint16)
//  2:   Dpad     (uint8: bit0=Up, bit1=Down, bit2=Left, bit3=Right)
//  3:   Triggers (uint8: bit6=ZL, bit7=ZR; other bits reserved)
//  4:   Stick    (uint8: bit0=Left, bit1=Right, bit2=Up, bit3=Down)
const (
	wireSeq     = 0
	wireDpad    = 2
	wireTriggers = 3
	wireStick   = 4
)

// Dpad bits in wire byte 2.
const (
	DpadUp    = 1 << 0
	DpadDown  = 1 << 1
	DpadLeft  = 1 << 2
	DpadRight = 1 << 3
)

// Trigger bits in wire byte 3.
const (
	TriggerZL = 1 << 6
	TriggerZR = 1 << 7
)

// Stick bits in wire byte 4.
const (
	StickLeft  = 1 << 0
	StickRight = 1 << 1
	StickUp    = 1 << 2
	StickDown  = 1 << 3
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
	var dpad, triggers uint8
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
		triggers |= TriggerZL
	}
	if s.TriggerRight {
		triggers |= TriggerZR
	}
	dst[wireDpad] = dpad
	dst[wireTriggers] = triggers
	var stick uint8
	if s.StickLeft {
		stick |= StickLeft
	}
	if s.StickRight {
		stick |= StickRight
	}
	if s.StickUp {
		stick |= StickUp
	}
	if s.StickDown {
		stick |= StickDown
	}
	dst[wireStick] = stick
	return dst[:PacketSize]
}

// Unmarshal decodes a PacketSize-byte packet into a new State. The wire structure is defined in this package.
func Unmarshal(data []byte) (*State, error) {
	if len(data) != PacketSize {
		return nil, errors.New("packet must be 5 bytes")
	}
	s := &State{}
	s.Sequence = binary.BigEndian.Uint16(data[wireSeq : wireSeq+2])
	dpad := data[wireDpad]
	s.DpadUp = dpad&DpadUp != 0
	s.DpadDown = dpad&DpadDown != 0
	s.DpadLeft = dpad&DpadLeft != 0
	s.DpadRight = dpad&DpadRight != 0
	trig := data[wireTriggers]
	s.TriggerLeft = trig&TriggerZL != 0
	s.TriggerRight = trig&TriggerZR != 0
	stk := data[wireStick]
	s.StickLeft = stk&StickLeft != 0
	s.StickRight = stk&StickRight != 0
	s.StickUp = stk&StickUp != 0
	s.StickDown = stk&StickDown != 0
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

func (s *State) String() string {
	bit := func(b bool) int { if b { return 1 }; return 0 }
	return fmt.Sprintf("DPad↑→↓←: %d%d%d%d, Stick↑→↓←: %d%d%d%d, Trig←→: %d%d, Seq: %d",
		bit(s.DpadUp), bit(s.DpadRight), bit(s.DpadDown), bit(s.DpadLeft),
		bit(s.StickUp), bit(s.StickRight), bit(s.StickDown), bit(s.StickLeft),
		bit(s.TriggerLeft), bit(s.TriggerRight), s.Sequence)
}
