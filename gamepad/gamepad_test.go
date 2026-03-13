package gamepad

import (
	"bytes"
	"testing"
)

func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	s := &State{
		Sequence:     42,
		DpadUp:       true,
		DpadRight:    true,
		StickLeft:    true,
		TriggerRight: true,
	}
	dst := make([]byte, PacketSize)
	out := s.Marshal(dst)
	if out == nil {
		t.Fatal("Marshal returned nil")
	}
	if len(out) != PacketSize {
		t.Fatalf("Marshal wrote %d bytes, want %d", len(out), PacketSize)
	}

	got, err := Unmarshal(out)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if got.Sequence != s.Sequence {
		t.Errorf("Sequence: got %d, want %d", got.Sequence, s.Sequence)
	}
	if got.DpadUp != s.DpadUp || got.DpadDown != s.DpadDown || got.DpadLeft != s.DpadLeft || got.DpadRight != s.DpadRight {
		t.Errorf("Dpad: got %v,%v,%v,%v want %v,%v,%v,%v",
			got.DpadUp, got.DpadDown, got.DpadLeft, got.DpadRight,
			s.DpadUp, s.DpadDown, s.DpadLeft, s.DpadRight)
	}
	if got.StickUp != s.StickUp || got.StickDown != s.StickDown || got.StickLeft != s.StickLeft || got.StickRight != s.StickRight {
		t.Errorf("Stick: got %v,%v,%v,%v want %v,%v,%v,%v",
			got.StickUp, got.StickDown, got.StickLeft, got.StickRight,
			s.StickUp, s.StickDown, s.StickLeft, s.StickRight)
	}
	if got.TriggerLeft != s.TriggerLeft || got.TriggerRight != s.TriggerRight {
		t.Errorf("Triggers: got L=%v R=%v want L=%v R=%v", got.TriggerLeft, got.TriggerRight, s.TriggerLeft, s.TriggerRight)
	}
}

func TestUnmarshalWrongSize(t *testing.T) {
	_, err := Unmarshal(make([]byte, 8))
	if err == nil {
		t.Error("Unmarshal(8 bytes) want error")
	}
	_, err = Unmarshal(make([]byte, 20))
	if err == nil {
		t.Error("Unmarshal(20 bytes) want error")
	}
}

func TestMarshalShortDst(t *testing.T) {
	s := &State{}
	if s.Marshal(make([]byte, 4)) != nil {
		t.Error("Marshal(short dst) should return nil")
	}
}

func TestMarshalSequenceBigEndian(t *testing.T) {
	s := &State{Sequence: 0x1234}
	dst := make([]byte, PacketSize)
	s.Marshal(dst)
	// Big-endian: 0x12, 0x34 in first two bytes
	if dst[0] != 0x12 || dst[1] != 0x34 {
		t.Errorf("Sequence bytes: got %02x %02x, want 12 34", dst[0], dst[1])
	}
}

func TestEmptyStateRoundTrip(t *testing.T) {
	s := &State{Sequence: 1}
	dst := make([]byte, PacketSize)
	s.Marshal(dst)
	got, err := Unmarshal(dst)
	if err != nil {
		t.Fatal(err)
	}
	if got.Sequence != 1 {
		t.Errorf("Sequence: got %d", got.Sequence)
	}
	if got.DpadUp || got.DpadDown || got.DpadLeft || got.DpadRight {
		t.Error("Dpad should be all false")
	}
	if got.StickUp || got.StickDown || got.StickLeft || got.StickRight {
		t.Error("Stick should be all false")
	}
	if got.TriggerLeft || got.TriggerRight {
		t.Error("Triggers should be false")
	}
	// Reserved bytes stay zero
	if !bytes.Equal(dst[8:12], []byte{0, 0, 0, 0}) {
		t.Errorf("Reserved bytes not zero: %v", dst[8:12])
	}
}
