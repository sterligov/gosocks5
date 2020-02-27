package socks5

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

// Socks server config
type Socks struct {
	Port           string
	MaxConnections int
	Deadline       time.Duration
	Authorization  Authorizer
}

const (
	statusOK = iota
	errorServer
	errorRule
	errorNet
	errorHost
	errorRefusing
	errorTimeoutTTL
	errorUnknownCommand
	errorAddressTypeNotSupported

	socksV5        = 0x05
	commandConnect = 0x01
	bindConnect    = 0x02
)

func (s *Socks) Run() {
	listener, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	guard := make(chan struct{}, s.MaxConnections)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		select {
		case guard <- struct{}{}:
		case <-time.After(100 * time.Millisecond):
			conn.Close()
			log.Println("too many connections")
			continue
		}

		go func() {
			err := s.handle(&connection{conn, s.Deadline})
			conn.Close()
			<-guard
			if err != nil {
				log.Println(err)
			}
		}()
	}
}

func (s *Socks) handle(conn *connection) error {
	buf := make([]byte, 1)
	_, err := conn.Read(buf)
	if err != nil {
		return err
	}
	if buf[0] != socksV5 {
		return fmt.Errorf("unsupported socks version %d", buf[0])
	}

	err = authorize(s.Authorization, conn)
	if err != nil {
		return err
	}

	buf = make([]byte, 3)
	_, err = conn.Read(buf)
	if err != nil {
		return err
	}
	if buf[0] != socksV5 {
		return fmt.Errorf("unsupported socks version %d", buf[0])
	}

	addr, err := destinationAddr(conn)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte{socksV5, statusOK, 0, ipv4, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return err
	}

	return dial(conn, addr)
}

func dial(client io.ReadWriter, raddr *net.TCPAddr) error {
	server, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return err
	}
	defer server.Close()

	errCh := make(chan error, 2)
	go copyData(server, client, errCh)
	go copyData(client, server, errCh)

	for i := 0; i < 2; i++ {
		err := <-errCh
		if err != nil {
			return err
		}
	}

	return nil
}

func copyData(dst io.Writer, src io.Reader, errCh chan<- error) {
	// buf := make([]byte, 32*1024)
	// for {
	// 	nr, er := src.Read(buf)
	// 	if nr > 0 {
	// 		nw, ew := dst.Write(buf[0:nr])
	// 		if ew != nil {
	// 			errCh <- ew
	// 			return
	// 		}
	// 		if nr != nw {
	// 			errCh <- io.ErrShortWrite
	// 			return
	// 		}
	// 	}
	// 	if er != nil {
	// 		if er != io.EOF {
	// 			errCh <- er
	// 			return
	// 		}
	// 		errCh <- nil
	// 		return
	// 	}
	// }
	_, err := io.Copy(dst, src)
	log.Println("EERROR10: ", err)
	errCh <- err
}
