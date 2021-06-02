package fquic

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/lucas-clemente/quic-go"
)

type Listener struct {
	lis quic.Listener
}

func newListener(lis quic.Listener) *Listener {
	return &Listener{
		lis: lis,
	}
}

func Listen(address string) (*Listener, error) {
	return new(ListenConfig).Listen(address)
}

func Server(conn net.PacketConn) (*Listener, error) {
	return new(ListenConfig).Server(conn)
}

func (lis *Listener) Close() error {
	return lis.lis.Close()
}

func (lis *Listener) Accept() (net.Conn, error) {
	conn, err := lis.AcceptQUIC(context.Background())
	if err != nil {
		return nil, err
	}

	return conn.AcceptStream(context.Background())
}

func (lis *Listener) AcceptQUIC(ctx context.Context) (*Conn, error) {
	s, err := lis.lis.Accept(ctx)
	if err != nil {
		return nil, err
	}

	return newConn(s), nil
}

func (lis *Listener) Addr() net.Addr {
	return lis.lis.Addr()
}

// ListenConfig defines configuration details for a Listener.
// Undocumented fields behave the same as identically named fields do
// in Dialer.
type ListenConfig struct {
	TLSConfig  *tls.Config
	QUICConfig *quic.Config
	Protocol   string
}

func (lc *ListenConfig) dialer() *Dialer {
	return &Dialer{
		TLSConfig:  lc.TLSConfig,
		QUICConfig: lc.QUICConfig,
		Protocol:   lc.Protocol,
	}
}

func (lc *ListenConfig) Listen(address string) (*Listener, error) {
	lis, err := quic.ListenAddr(address, lc.dialer().tlsConfig(), lc.QUICConfig)
	if err != nil {
		return nil, err
	}
	return newListener(lis), nil
}

func (lc *ListenConfig) Server(conn net.PacketConn) (*Listener, error) {
	lis, err := quic.Listen(conn, lc.dialer().tlsConfig(), lc.QUICConfig)
	if err != nil {
		return nil, err
	}
	return newListener(lis), nil
}
