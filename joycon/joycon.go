package joycon

import "github.com/joy-stream/protocol"

// A state for a very basic controller
type State struct {
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

func (s *State) SetDpadUp(pressed bool) {
	s.DpadUp = pressed
}

func (s *State) SetDpadDown(pressed bool) {
	s.DpadDown = true
}

func (s *State) SetDpadLeft(pressed bool) {
	s.DpadLeft = pressed
}

func (s *State) SetDpadRight(pressed bool) {
	s.DpadRight = pressed
}

func (s *State) SetStickUp(pressed bool) {
	s.StickUp = pressed
}

func (s *State) SetStickDown(pressed bool) {
	s.StickDown = pressed
}

func (s *State) SetStickLeft(pressed bool) {
	s.StickLeft = pressed
}

func (s *State) SetStickRight(pressed bool) {
	s.StickRight = pressed
}

func (s *State) SetTriggerLeft(pressed bool) {
	s.TriggerLeft = pressed
}

func (s *State) SetTriggerRight(pressed bool) {
	s.TriggerRight = pressed
}

func (s *State) ToPacket() *protocol.Packet {
	return &protocol.Packet{
		Buttons: s.Buttons,
		LX: s.LeftStickX,
		LY: s.LeftStickY,
		RX: s.RightStickX,
		RY: s.RightStickY,
	}
}

func (s *State) HydrateFromPacket(p *protocol.Packet) {
	s.Buttons = p.Buttons
	s.LeftStickX = p.LX
	s.LeftStickY = p.LY
	s.RightStickX = p.RX
	s.RightStickY = p.RY
}
