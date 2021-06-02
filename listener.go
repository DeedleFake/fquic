package fquic

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"

	"github.com/lucas-clemente/quic-go"
)

// Listener listens for incoming connections.
type Listener struct {
	lis quic.Listener
}

// Listen creates a Listener that listens for connections on the given
// local address.
func Listen(address string) (*Listener, error) {
	return new(ListenConfig).Listen(address)
}

// Server creates a Listener that listens for connections using the
// given net.PacketConn.
func Server(conn net.PacketConn) (*Listener, error) {
	return new(ListenConfig).Server(conn)
}

// Listener returns the underlying quic-go Listener.
func (lis *Listener) Listener() quic.Listener {
	return lis.lis
}

// NetListener returns a wrapper around the Listener that implements
// net.Listener. In order to do this, it immediately creates a new
// Stream after a connection is accepted. If outgoing is false, the
// stream is expected to be initiated by the remote end, whereas if
// it is true it is created at the local end of the connection. If
// outgoing is true, unidirectional specifies the type of stream
// created, but otherwise it is ignored.
func (lis *Listener) NetListener(outgoing, unidirectional bool) net.Listener {
	return netListener{
		Listener: lis,

		outgoing:       outgoing,
		unidirectional: unidirectional,
	}
}

// Close closes the listener. This will also close all connections
// that were accepted by this listener.
func (lis *Listener) Close() error {
	return lis.lis.Close()
}

// Accept accepts a new connection.
// Note that because Conn does not implement net.Conn, Listener does
// not implement net.Listener. For a simple workaround, see the
// NetListener method.
func (lis *Listener) Accept(ctx context.Context) (*Conn, error) {
	s, err := lis.lis.Accept(ctx)
	if err != nil {
		return nil, err
	}

	return newConn(s), nil
}

// Addr returns the local address that the Listener is listening on.
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

	// If TLSConfig contains no certificates, GetCertificate is used to
	// get one as the Listener must have a certificate. If
	// GetCertificate is nil, GenerateCert is used.
	GetCertificate func() (tls.Certificate, error)
}

func (lc *ListenConfig) dialer() *Dialer {
	return &Dialer{
		TLSConfig:  lc.TLSConfig,
		QUICConfig: lc.QUICConfig,
		Protocol:   lc.Protocol,
	}
}

func (lc *ListenConfig) tlsConfig() (*tls.Config, error) {
	conf := lc.dialer().tlsConfig()

	if len(conf.Certificates) == 0 {
		get := lc.GetCertificate
		if get == nil {
			get = GenerateCert
		}

		cert, err := get()
		if err != nil {
			return nil, fmt.Errorf("get certificate: %w", err)
		}

		conf.Certificates = []tls.Certificate{cert}
	}

	return conf, nil
}

// Listen listens for incoming connections on the specified address.
func (lc *ListenConfig) Listen(address string) (*Listener, error) {
	tlsConfig, err := lc.tlsConfig()
	if err != nil {
		return nil, err
	}

	lis, err := quic.ListenAddr(address, tlsConfig, lc.QUICConfig)
	if err != nil {
		return nil, err
	}

	return &Listener{lis: lis}, nil
}

// Server listens for incoming connections using the provided
// net.PacketConn.
func (lc *ListenConfig) Server(conn net.PacketConn) (*Listener, error) {
	tlsConfig, err := lc.tlsConfig()
	if err != nil {
		return nil, err
	}

	lis, err := quic.Listen(conn, tlsConfig, lc.QUICConfig)
	if err != nil {
		return nil, err
	}

	return &Listener{lis: lis}, nil
}

type netListener struct {
	*Listener

	outgoing       bool
	unidirectional bool
}

func (lis netListener) Accept() (net.Conn, error) {
	conn, err := lis.Listener.Accept(context.Background())
	if err != nil {
		return nil, err
	}

	if lis.outgoing {
		return conn.NewStream(lis.unidirectional)
	}
	return conn.AcceptStream(context.Background())
}

// GenerateCert generates a new, self-signed certificate.
func GenerateCert() (tls.Certificate, error) {
	pubkey, privkey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
	}
	xcert, err := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template,
		pubkey,
		privkey,
	)
	if err != nil {
		return tls.Certificate{}, err
	}

	xkey, err := x509.MarshalPKCS8PrivateKey(privkey)
	if err != nil {
		return tls.Certificate{}, err
	}

	pkey := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: xkey,
	})
	pcert := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: xcert,
	})

	return tls.X509KeyPair(pcert, pkey)
}
