package quic

import (
	"context"
	"crypto/tls"
	"net"
)

type Conn struct {
	conn  net.PacketConn
	raddr net.Addr
}

func Dial(network, address string) (*Conn, error) {
	return new(Dialer).Dial(network, address)
}

func Client(conn net.PacketConn, raddr net.Addr) (*Conn, error) {
	panic("Not implemented.")
}

func (c *Conn) Close() error {
	panic("Not implemented.")
}

func (c *Conn) AcceptStream() (*Stream, error) {
	panic("Not implemented.")
}

func (c *Conn) NewStream(unidirectional bool) (*Stream, error) {
	panic("Not implemented.")
}

func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.RemoteAddr()
}

type Dialer struct {
	TLSConfig *tls.Config
}

func (d *Dialer) Dial(network, address string) (*Conn, error) {
	return d.DialContext(context.Background(), network, address)
}

func (d *Dialer) DialContext(ctx context.Context, network, address string) (*Conn, error) {
	panic("Not implemented.")
}
