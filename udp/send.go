package udp

import "net"

// Sender sends UDP packets to a fixed server address using a single connection.
type Sender struct {
	conn *net.UDPConn
	addr *net.UDPAddr
}

// NewSender resolves serverAddr (e.g. "127.0.0.1:7355") and opens a UDP connection.
func NewSender(serverAddr string) (*Sender, error) {
	addr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	return &Sender{conn: conn, addr: addr}, nil
}

// Send writes data to the server. Caller passes the full packet (e.g. marshaled gamepad state).
func (s *Sender) Send(data []byte) error {
	_, err := s.conn.Write(data)
	return err
}

// Close closes the UDP connection.
func (s *Sender) Close() error {
	if s.conn == nil {
		return nil
	}
	return s.conn.Close()
}
