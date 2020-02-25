package socks5

import (
	"fmt"
	"io"
	"net"
)

const (
	ipv4       = 0x01
	domainName = 0x03
	ipv6       = 0x04
)

func destinationAddr(rw io.ReadWriter) (*net.TCPAddr, error) {
	buf := make([]byte, 1)
	_, err := rw.Read(buf)
	if err != nil {
		return nil, err
	}

	var ip net.IP
	if buf[0] == ipv4 {
		ip = make([]byte, 4)
		_, err := rw.Read(ip)
		if err != nil {
			return nil, err
		}
	} else if buf[0] == domainName {
		buf = make([]byte, 1)
		_, err := rw.Read(buf)
		if err != nil {
			return nil, err
		}

		domain := make([]byte, buf[0])
		_, err = rw.Read(domain)
		if err != nil {
			return nil, err
		}

		ips, err := net.LookupIP(string(domain))
		if err != nil {
			return nil, err
		}
		ip = ips[0]
	} else if buf[0] == ipv6 {
		ip = make([]byte, 16)
		_, err := rw.Read(ip)
		if err != nil {
			return nil, err
		}
	} else {
		rw.Write([]byte{socksV5, errorAddressTypeNotSupported})
		return nil, fmt.Errorf("unknown address type")
	}

	buf = make([]byte, 2)
	_, err = rw.Read(buf)
	if err != nil {
		return nil, err
	}
	port := int(buf[0])<<8 + int(buf[1])

	return &net.TCPAddr{IP: ip, Port: port}, nil
}
