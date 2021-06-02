package fquic

import "net"

type Listener struct {
	conn net.PacketConn
}

func Listen(network, address string) (*Listener, error) {
	panic("Not implemented.")
}

func Server(conn net.PacketConn) (*Listener, error) {
	panic("Not implemented.")
}

func (lis *Listener) Close() error {
	panic("Not implemented.")
}

func (lis *Listener) Accept() (net.Conn, error) {
	conn, err := lis.AcceptQUIC()
	if err != nil {
		return nil, err
	}

	return conn.AcceptStream()
}

func (lis *Listener) AcceptQUIC() (*Conn, error) {
	panic("Not implemented.")
}

func (lis *Listener) Addr() net.Addr {
	return lis.conn.LocalAddr()
}
