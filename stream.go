package fquic

import (
	"errors"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go"
)

var (
	ErrWriteOnly = errors.New("write-only stream")
	ErrReadOnly  = errors.New("read-only stream")
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

func (s *Stream) Stream() quic.Stream {
	qs, _ := s.s.(quic.Stream)
	return qs
}

func (s *Stream) ReceiveStream() quic.ReceiveStream {
	return s.r
}

func (s *Stream) SendStream() quic.SendStream {
	return s.s
}

func (s *Stream) CanReceive() bool {
	return s.r != nil
}

func (s *Stream) CanSend() bool {
	return s.s != nil
}

func (s *Stream) Read(buf []byte) (int, error) {
	if s.r == nil {
		return 0, ErrWriteOnly
	}

	return s.r.Read(buf)
}

func (s *Stream) Write(data []byte) (int, error) {
	if s.s == nil {
		return 0, ErrReadOnly
	}

	return s.s.Write(data)
}

func (s *Stream) Close() error {
	if s.s == nil {
		return nil
	}
	return s.s.Close()
}

func (s *Stream) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}

func (s *Stream) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *Stream) SetDeadline(t time.Time) error {
	if s.s == nil {
		return ErrReadOnly
	}
	if s.r == nil {
		return ErrWriteOnly
	}

	type deadliner interface {
		SetDeadline(time.Time) error
	}
	return s.s.(deadliner).SetDeadline(t)
}

func (s *Stream) SetReadDeadline(t time.Time) error {
	if s.r == nil {
		return ErrWriteOnly
	}
	return s.r.SetReadDeadline(t)
}

func (s *Stream) SetWriteDeadline(t time.Time) error {
	if s.s == nil {
		return ErrReadOnly
	}
	return s.s.SetWriteDeadline(t)
}
