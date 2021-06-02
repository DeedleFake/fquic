package fquic

import (
	"net"
	"time"
)

type Stream struct {
	conn *Conn
}

func (s *Stream) Conn() *Conn {
	return s.conn
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
