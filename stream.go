package fquic

import (
	"errors"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go"
)

// These errors are panicked by Stream methods that are called on
// incorrectly used unidirectional streams.
var (
	ErrWriteOnly = errors.New("write-only stream")
	ErrReadOnly  = errors.New("read-only stream")
)

// A Stream is an open QUIC stream inside of a connection.
//
// A Stream can be bidirectional, read-only, or write-only. Calling
// read-related methods on a write-only stream and vice versa will
// result in a panic.
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

// Conn returns the connection associated with the Stream.
func (s *Stream) Conn() *Conn {
	return s.conn
}

// Stream returns the underlying bidirectional quic-go Stream. If the
// Stream is unidirectional in either direction, this returns nil.
func (s *Stream) Stream() quic.Stream {
	qs, _ := s.s.(quic.Stream)
	return qs
}

// ReceiveStream returns the underlying quic-go ReceiveStream, or nil
// if the stream is write-only.
func (s *Stream) ReceiveStream() quic.ReceiveStream {
	return s.r
}

// SendStream returns the underlying quic-go SendStream, or nil if the
// stream is read-only.
func (s *Stream) SendStream() quic.SendStream {
	return s.s
}

// CanRead returns true if the stream can be read from.
func (s *Stream) CanRead() bool {
	return s.r != nil
}

// CanWrite returns true if the stream can be written to.
func (s *Stream) CanWrite() bool {
	return s.s != nil
}

func (s *Stream) Read(buf []byte) (int, error) {
	if s.r == nil {
		panic(ErrWriteOnly)
	}

	return s.r.Read(buf)
}

func (s *Stream) Write(data []byte) (int, error) {
	if s.s == nil {
		panic(ErrReadOnly)
	}

	return s.s.Write(data)
}

// Close closes a writable stream. Unlike other direction-dependant
// methods methods, if this stream is read-only, this is a no-op.
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

// SetDeadline implements net.Conn's SetDeadline method. Note that
// this method only works on bidirectional streams.
//
// TODO: Make this use the available unidirectional method
// automatically if the stream is unidirectional?
func (s *Stream) SetDeadline(t time.Time) error {
	if s.s == nil {
		panic(ErrReadOnly)
	}
	if s.r == nil {
		panic(ErrWriteOnly)
	}

	type deadliner interface {
		SetDeadline(time.Time) error
	}
	return s.s.(deadliner).SetDeadline(t)
}

func (s *Stream) SetReadDeadline(t time.Time) error {
	if s.r == nil {
		panic(ErrWriteOnly)
	}
	return s.r.SetReadDeadline(t)
}

func (s *Stream) SetWriteDeadline(t time.Time) error {
	if s.s == nil {
		panic(ErrReadOnly)
	}
	return s.s.SetWriteDeadline(t)
}
