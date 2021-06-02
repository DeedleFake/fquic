package fquic

import (
	"net"
	"time"

	"github.com/lucas-clemente/quic-go"
)

type Stream struct {
	conn *Conn
	s    quic.SendStream
	r    quic.ReceiveStream
}

func newStream(conn *Conn, s quic.SendStream, r quic.ReceiveStream) *Stream {
	return &Stream{
		conn: conn,
		s:    s,
		r:    r,
	}
}

func (s *Stream) Conn() *Conn {
	return s.conn
}

func (s *Stream) CanSend() bool {
	return s.s != nil
}

func (s *Stream) CanReceive() bool {
	return s.r != nil
}

func (s *Stream) Read(buf []byte) (int, error) {
	panic("Not implemented.")
}

func (s *Stream) Write(data []byte) (int, error) {
	panic("Not implemented.")
}

func (s *Stream) Close() error {
	panic("Not implemented.")
}

func (s *Stream) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}

func (s *Stream) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *Stream) SetDeadline(t time.Time) error {
	panic("Not implemented.")
}

func (s *Stream) SetReadDeadline(t time.Time) error {
	panic("Not implemented.")
}

func (s *Stream) SetWriteDeadline(t time.Time) error {
	panic("Not implemented.")
}
