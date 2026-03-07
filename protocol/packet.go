package protocol

import (
	"encoding/binary"
	"errors"
)

// Packet size in bytes.
const PacketSize = 12

// Button bitmask bits (uint16).
const (
	ButtonA       = 1 << 0
	ButtonB       = 1 << 1
	ButtonX       = 1 << 2
	ButtonY       = 1 << 3
	ButtonL       = 1 << 4
	ButtonR       = 1 << 5
	ButtonZL      = 1 << 6
	ButtonZR      = 1 << 7
	ButtonPlus    = 1 << 8
	ButtonMinus   = 1 << 9
	ButtonLClick  = 1 << 10
	ButtonRClick  = 1 << 11
	ButtonHome    = 1 << 12
	ButtonCapture = 1 << 13
)

// Packet represents a decoded 12-byte controller state packet.
// Byte 0-1: sequence (uint16 big endian)
// Byte 2-3: button bitmask (uint16)
// Byte 4: LX, 5: LY, 6: RX, 7: RY (0-255, 128=center)
// Byte 8: dpad, 9: misc, 10-11: reserved
type Packet struct {
	Sequence uint16
	Buttons  uint16
	LX       uint8
	LY       uint8
	RX       uint8
	RY       uint8
	Dpad     uint8
	Misc     uint8
}

// ParsePacket decodes a 12-byte UDP payload into a Packet.
// Returns an error if len(data) != 12.
func ParsePacket(data []byte) (*Packet, error) {
	if len(data) != PacketSize {
		return nil, errors.New("packet must be 12 bytes")
	}
	p := &Packet{
		Sequence: binary.BigEndian.Uint16(data[0:2]),
		Buttons:  binary.BigEndian.Uint16(data[2:4]),
		LX:       data[4],
		LY:       data[5],
		RX:       data[6],
		RY:       data[7],
		Dpad:     data[8],
		Misc:     data[9],
	}
	return p, nil
}

// Marshal writes the packet to a 12-byte slice (big-endian). Returns nil if dst is too short.
func (p *Packet) Marshal(dst []byte) []byte {
	if len(dst) < PacketSize {
		return nil
	}
	binary.BigEndian.PutUint16(dst[0:2], p.Sequence)
	binary.BigEndian.PutUint16(dst[2:4], p.Buttons)
	dst[4], dst[5], dst[6], dst[7] = p.LX, p.LY, p.RX, p.RY
	dst[8], dst[9] = p.Dpad, p.Misc
	dst[10], dst[11] = 0, 0
	return dst[:PacketSize]
}

// AxisToFloat converts a 0-255 axis value to -1..1 (128 = center).
func AxisToFloat(v uint8) float64 {
	return float64(int8(v-128)) / 127.0
}

// ButtonNames returns a slice of pressed button names for display.
func (p *Packet) ButtonNames() []string {
	var names []string
	if p.Buttons&ButtonA != 0 {
		names = append(names, "A")
	}
	if p.Buttons&ButtonB != 0 {
		names = append(names, "B")
	}
	if p.Buttons&ButtonX != 0 {
		names = append(names, "X")
	}
	if p.Buttons&ButtonY != 0 {
		names = append(names, "Y")
	}
	if p.Buttons&ButtonL != 0 {
		names = append(names, "L")
	}
	if p.Buttons&ButtonR != 0 {
		names = append(names, "R")
	}
	if p.Buttons&ButtonZL != 0 {
		names = append(names, "ZL")
	}
	if p.Buttons&ButtonZR != 0 {
		names = append(names, "ZR")
	}
	if p.Buttons&ButtonPlus != 0 {
		names = append(names, "+")
	}
	if p.Buttons&ButtonMinus != 0 {
		names = append(names, "-")
	}
	if p.Buttons&ButtonLClick != 0 {
		names = append(names, "LCLICK")
	}
	if p.Buttons&ButtonRClick != 0 {
		names = append(names, "RCLICK")
	}
	if p.Buttons&ButtonHome != 0 {
		names = append(names, "HOME")
	}
	if p.Buttons&ButtonCapture != 0 {
		names = append(names, "CAPTURE")
	}
	return names
}

// Dpad bits (same as client).
const (
	DpadUp    = 1 << 0
	DpadDown  = 1 << 1
	DpadLeft  = 1 << 2
	DpadRight = 1 << 3
)

// DpadNames returns D-pad direction names for display (e.g. "Haut", "Gauche").
func (p *Packet) DpadNames() []string {
	var names []string
	if p.Dpad&DpadUp != 0 {
		names = append(names, "Haut")
	}
	if p.Dpad&DpadDown != 0 {
		names = append(names, "Bas")
	}
	if p.Dpad&DpadLeft != 0 {
		names = append(names, "Gauche")
	}
	if p.Dpad&DpadRight != 0 {
		names = append(names, "Droite")
	}
	return names
}
