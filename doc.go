// Package fquic provides a friendly wrapper around quic-go.
//
// Echo server:
//    func server(address string) {
//    	lis, err := fquic.Listen("quic-echo-example", address)
//    	if err != nil {
//    		log.Fatalf("open listener: %v", err)
//    	}
//    	defer lis.Close()
//
//    	for {
//    		c, err := lis.Accept(context.Background())
//    		if err != nil {
//    			log.Fatalf("accept connection: %v", err)
//    		}
//
//    		go func() {
//    			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//    			defer cancel()
//
//    			s, err := c.AcceptStream(ctx)
//    			if err != nil {
//    				log.Printf("accept stream: %v", err)
//    				return
//    			}
//    			defer s.Close()
//    			if !s.CanWrite() {
//    				log.Printf("got read-only stream")
//    				return
//    			}
//
//    			var buf strings.Builder
//    			r := io.TeeReader(s, &buf)
//    			_, err = io.Copy(s, r)
//    			if err != nil {
//    				var apperr *quic.ApplicationError
//    				if errors.As(err, &apperr) {
//    					log.Printf("disconnect %v: %q", apperr.ErrorCode, apperr.ErrorMessage)
//    					log.Printf("echoed: %q", buf.String())
//    					return
//    				}
//
//    				log.Printf("echo error: %v", err)
//    				return
//    			}
//    		}()
//    	}
//    }
//
// Echo client:
//    func client(address string, data []byte) {
//    	d := fquic.Dialer{
//    		TLSConfig: &tls.Config{
//    			InsecureSkipVerify: true,
//    		},
//    		Protocol: "quic-echo-example",
//    	}
//    	c, err := d.Dial(address)
//    	if err != nil {
//    		log.Fatalf("dial: %v", err)
//    	}
//    	defer c.Close()
//
//    	s, err := c.NewStream(false)
//    	if err != nil {
//    		log.Fatalf("open stream: %v", err)
//    	}
//    	defer s.Close()
//
//    	_, err = s.Write(data)
//    	if err != nil {
//    		log.Fatalf("write data: %v", err)
//    	}
//
//    	buf := make([]byte, len(data))
//    	_, err = io.ReadFull(s, buf)
//    	if err != nil {
//    		log.Fatalf("read data: %v", err)
//    	}
//
//    	fmt.Printf("Got echo: %q\n", buf)
//    }
package fquic
