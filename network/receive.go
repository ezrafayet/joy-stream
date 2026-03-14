package network

import (
	"errors"
	"net"
)

// PacketHandler is called for each received UDP packet. data is the payload;
// the slice is only valid for the duration of the call—copy it if you need to keep it.
type PacketHandler func(data []byte, from net.Addr)

// Receiver listens on a UDP address and invokes a callback for each packet.
type Receiver struct {
	conn *net.UDPConn
}

// NewReceiver listens on listenAddr (e.g. ":7355"). Call Start to begin receiving.
func NewReceiver(listenAddr string) (*Receiver, error) {
	conn, err := net.ListenPacket("udp", listenAddr)
	if err != nil {
		return nil, err
	}
	udpConn, ok := conn.(*net.UDPConn)
	if !ok {
		conn.Close()
		return nil, errors.New("udp listen: unexpected connection type")
	}
	return &Receiver{conn: udpConn}, nil
}

// Start runs the receive loop in a goroutine, calling handler for each packet.
// Returns immediately. Call Close to stop receiving.
func (r *Receiver) Start(handler PacketHandler) {
	if r.conn == nil || handler == nil {
		return
	}
	go func() {
		buf := make([]byte, 65535)
		for {
			n, from, err := r.conn.ReadFrom(buf)
			if err != nil {
				return
			}
			if n == 0 {
				continue
			}
			handler(buf[:n:n], from)
		}
	}()
}

// Close stops the receiver and closes the UDP connection.
func (r *Receiver) Close() error {
	if r.conn == nil {
		return nil
	}
	return r.conn.Close()
}
