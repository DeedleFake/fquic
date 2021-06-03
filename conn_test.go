package fquic_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"sync"
	"testing"

	"github.com/DeedleFake/fquic"
	"github.com/lucas-clemente/quic-go"
	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	testData := []byte("some test data that is all lowercase")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	t.Cleanup(cancel)

	var wg sync.WaitGroup
	wg.Add(2)

	addr := make(chan net.Addr, 1)
	go func() {
		defer wg.Done()

		lis, err := fquic.Listen("test-simple", "localhost:0")
		assert.NoError(t, err, "open listener")
		addr <- lis.Addr()
		defer func() {
			err := lis.Close()
			assert.NoError(t, err, "close listener")
		}()

		c, err := lis.Accept(ctx)
		assert.NoError(t, err, "accept connection")
		defer func() {
			err := c.Close()
			assert.NoError(t, err, "close connection")
		}()

		s, err := c.AcceptStream(ctx)
		assert.NoError(t, err, "accept stream")
		defer func() {
			err := s.Close()
			assert.NoError(t, err, "close stream")
		}()
		assert.True(t, s.CanWrite(), "bidirectional stream")

		var buf [1024]byte
		n, err := s.Read(buf[:])
		if err != io.EOF {
			assert.NoError(t, err, "read")
		}
		assert.Equal(t, testData, buf[:n], "data")

		n, err = s.Write(bytes.ToUpper(buf[:n]))
		assert.NoError(t, err, "write")
		assert.Equal(t, len(testData), n, "amount written")

		n, err = s.Read(make([]byte, 1))
		if err != nil {
			var apperr *quic.ApplicationError
			if errors.As(err, &apperr) {
				if apperr.ErrorCode == quic.ApplicationErrorCode(quic.NoError) {
					return
				}
			}
			assert.NoError(t, err, "wait for close")
		}
	}()

	go func() {
		defer wg.Done()

		d := &fquic.Dialer{
			TLSConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			Protocol: "test-simple",
		}
		c, err := d.Dial((<-addr).String())
		assert.NoError(t, err, "dial")
		defer func() {
			err := c.Close()
			assert.NoError(t, err, "close connection")
		}()

		s, err := c.NewStream(false)
		assert.NoError(t, err, "open stream")
		defer func() {
			err := s.Close()
			assert.NoError(t, err, "close stream")
		}()

		n, err := s.Write(testData)
		assert.NoError(t, err, "write")
		assert.Equal(t, len(testData), n, "amount written")

		var buf [1024]byte
		n, err = s.Read(buf[:])
		if err != io.EOF {
			assert.NoError(t, err, "read")
		}
		assert.Equal(t, bytes.ToUpper(testData), buf[:n], "data")
	}()

	wg.Wait()
}
