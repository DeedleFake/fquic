// Package fquic provides a friendly wrapper around quic-go.
//
// Echo server:
//    lc := fquic.ListenConfig{
//      TLSConfig: &tls.Config{
//        InsecureSkipVerify: true,
//      },
//      Protocol: "quic-echo-example",
//    }
//    lis, err := lc.Listen(address)
//    if err != nil {
//      log.Fatalf("open listener: %v", err)
//    }
//    defer lis.Close()
//
//    for {
//      c, err := lis.Accept(context.Background())
//      if err != nil {
//        log.Fatalf("accept connection: %v", err)
//      }
//
//      go func() {
//        ctx, cancel := context.WithTimeout(10 * time.Second)
//        defer cancel()
//
//        s, err := c.AcceptStream(ctx)
//        if err != nil {
//          log.Printf("accept stream: %v", err)
//          return
//        }
//        defer s.Close()
//        if !s.CanWrite() {
//          log.Printf("got read-only stream")
//          return
//        }
//        _, err = io.Copy(s, s)
//        if err != nil {
//          log.Printf("echo: %v", err)
//          return
//        }
//      }()
//    }
//
// Echo client:
//    d := fquic.Dialer{
//      TLSConfig: &tls.Config{
//        InsecureSkipVerify: true,
//      },
//      Protocol: "quic-echo-example",
//    }
//    c, err := d.Dial(address)
//    if err != nil {
//      log.Fatalf("dial: %v", err)
//    }
//    defer c.Close()
//
//    s, err := c.NewStream(true)
//    if err != nil {
//      log.Fatalf("open stream: %v", err)
//    }
//    defer s.Close()
//
//    _, err = io.WriteString(s, data)
//    if err != nil {
//      log.Fatalf("write data: %v", err)
//    }
//
//    buf := make([]byte, len(data))
//    _, err := io.ReadFull(s, buf)
//    if err != nil {
//      log.Fatalf("read data: %v", err)
//    }
//
//    fmt.Printf("Got echo: %q\n", buf)
package fquic
