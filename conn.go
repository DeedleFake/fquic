package fquic

import (
	"context"
	"crypto/tls"
	"net"
	"sync"

	"github.com/lucas-clemente/quic-go"
	"golang.org/x/sync/errgroup"
)

const (
	DefaultCloseCode    = 0
	DefaultCloseMessage = "closed"
)

// Conn represents an open QUIC connection.
//
// Note that Conn does not implement net.Conn, as QUIC connections
// require streams to be opened for actual data transfer.
type Conn struct {
	session quic.Session

	streams    chan *Stream
	streamErr  error
	streamLock sync.RWMutex

	closer sync.Once
	done   chan struct{}
}

func newConn(session quic.Session) *Conn {
	c := Conn{
		session: session,
		streams: make(chan *Stream),
		done:    make(chan struct{}),
	}
	go c.acceptStreams()
	return &c
}

// Dial connects to the specified address using protocol as the
// NextProtos specification of the TLS configuration.
func Dial(protocol, address string) (*Conn, error) {
	return (&Dialer{
		Protocol: protocol,
	}).Dial(address)
}

// Client connects to the remote address using the provided net.PacketConn. The
// host parameter is used for SNI. It uses protocol as the NextProtos
// specification of the TLS configuration.
func Client(protocol string, conn net.PacketConn, raddr net.Addr, host string) (*Conn, error) {
	return (&Dialer{
		Protocol: protocol,
	}).Client(conn, raddr, host)
}

func (c *Conn) acceptStreams() {
	defer close(c.streams)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-c.done
		cancel()
	}()

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		for {
			s, err := c.session.AcceptStream(ctx)
			if err != nil {
				return err
			}

			c.streams <- newStream(c, s, s)
		}
	})

	eg.Go(func() error {
		for {
			s, err := c.session.AcceptUniStream(ctx)
			if err != nil {
				return err
			}

			c.streams <- newStream(c, nil, s)
		}
	})

	err := eg.Wait()
	if err != nil {
		c.streamLock.Lock()
		defer c.streamLock.Unlock()

		c.streamErr = err
	}
}

// Close closes the connection using DefaultCloseCode and
// DefaultCloseMessage.
func (c *Conn) Close() error {
	return c.CloseWithError(DefaultCloseCode, DefaultCloseMessage)
}

// CloseWithError closes the connection with provided error code and
// message. Error codes are application-defined.
func (c *Conn) CloseWithError(code uint64, message string) error {
	c.closer.Do(func() {
		close(c.done)
	})
	return c.session.CloseWithError(quic.ApplicationErrorCode(code), message)
}

// AcceptStream accepts a stream initiated by the peer in a similar
// manner to a Listener accepting a connection. Note that the returned
// stream can be either bidirectional or read-only.
func (c *Conn) AcceptStream(ctx context.Context) (*Stream, error) {
	// TODO: Make sure that this returns the correct errors in different
	// types of situations.

	select {
	case <-ctx.Done():
		return nil, ctx.Err()

	case s, ok := <-c.streams:
		if !ok {
			c.streamLock.RLock()
			defer c.streamLock.RUnlock()

			return nil, c.streamErr
		}
		return s, nil
	}
}

// NewStream opens a new stream that can be written to. If
// unidirectional if false, the peer may also send data.
func (c *Conn) NewStream(unidirectional bool) (*Stream, error) {
	if unidirectional {
		s, err := c.session.OpenUniStream()
		if err != nil {
			return nil, err
		}
		return newStream(c, s, nil), nil
	}

	s, err := c.session.OpenStream()
	if err != nil {
		return nil, err
	}
	return newStream(c, s, s), nil
}

// LocalAddr returns the address of the local end of the connection.
func (c *Conn) LocalAddr() net.Addr {
	return c.session.LocalAddr()
}

// RemoteAddr returns the address of the remote end of the connection.
func (c *Conn) RemoteAddr() net.Addr {
	return c.session.RemoteAddr()
}

// Session returns the underlying quic-go Session. Be careful when
// using this.
func (c *Conn) Session() quic.Session {
	return c.session
}

// ReadDatagram reads a datagram packet from the connection, if
// datagram support is enabled by both the local machine and the peer.
//
// To enable datagram support, set EnableDatagrams to true in the
// quic.Config via either Dialer or ListenConfig.
func (c *Conn) ReadDatagram(buf []byte) (int, error) {
	return c.session.Read(buf)
}

// WriteDatagram writes a datagram packet to the connection, if
// datagram support is enabled by both the local machine and the peer.
//
// To enable datagram support, set EnableDatagrams to true in the
// quic.Config via either Dialer or ListenConfig.
func (c *Conn) WriteDatagram(data []byte) (int, error) {
	return c.session.Write(data)
}

// SupportsDatagrams returns true if both ends of the connection have
// datagram support.
func (c *Conn) SupportsDatagrams() bool {
	return c.session.ConnectionState().SupportsDatagrams
}

// A Dialer contains options for connecting to an address. Use of
// methods on a zero-value Dialer is equivalent to calling the
// similarly-named top-level functions in this package.
//
// It is safe to call methods on Dialer concurrently, but they should
// not be called concurrently to modifying fields.
type Dialer struct {
	// TLSConfig is the TLS configuration to use when dialing a new
	// connection. If it is a nil, a sane default configuration is used.
	TLSConfig *tls.Config

	// QUICConfig is the quic-go configuration to use when dialing a new
	// connection. If it is nil, a sane default configuration is used.
	QUICConfig *quic.Config

	// Protocol, if non-empty, is used to build the NextProtos
	// specification of TLSConfig. One or the other must be specified.
	// If neither are specified, dialing operations will panic. If both
	// are specified, Protocol will be prepended to the list specified
	// in NextProtos.
	Protocol string
}

func (d *Dialer) tlsConfig() *tls.Config {
	conf := d.TLSConfig.Clone()
	if conf == nil {
		conf = new(tls.Config)
	}

	if d.Protocol != "" {
		conf.NextProtos = append([]string{d.Protocol}, conf.NextProtos...)
	}

	if len(conf.NextProtos) == 0 {
		panic("no protocol specified")
	}

	return conf
}

// Dial connects to the specified address.
func (d *Dialer) Dial(address string) (*Conn, error) {
	return d.DialContext(context.Background(), address)
}

// DialContext connects to the specified address using the provided context.
func (d *Dialer) DialContext(ctx context.Context, address string) (*Conn, error) {
	session, err := quic.DialAddrContext(ctx, address, d.tlsConfig(), d.QUICConfig)
	if err != nil {
		return nil, err
	}

	return newConn(session), nil
}

// Client connects to the remote address using the provided net.PacketConn. The
// host parameter is used for SNI.
func (d *Dialer) Client(conn net.PacketConn, raddr net.Addr, host string) (*Conn, error) {
	return d.ClientContext(context.Background(), conn, raddr, host)
}

// ClientContext connects to the remote address using the provided
// net.PacketConn and context. The host parameter is used for SNI.
func (d *Dialer) ClientContext(ctx context.Context, conn net.PacketConn, raddr net.Addr, host string) (*Conn, error) {
	session, err := quic.DialContext(ctx, conn, raddr, host, d.tlsConfig(), d.QUICConfig)
	if err != nil {
		return nil, err
	}

	return newConn(session), nil
}
