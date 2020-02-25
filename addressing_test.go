package socks5

import (
	"bytes"
	"testing"
)

func TestDestinationAddrIP4(t *testing.T) {
	b := []byte{ipv4, 127, 0, 0, 1, 0, 80}
	buf := bytes.NewBuffer(b)
	addr, err := destinationAddr(buf)
	if err != nil {
		t.Errorf("Expected nil, got [%v]", err)
	}
	if addr.String() != "127.0.0.1:80" && addr.String() != "[::1]:80" {
		t.Errorf("Expected addr 127.0.0.1:80 or [::1]:80, got %s", addr.String())
	}
}

func TestDestinationAddrIP6(t *testing.T) {
	b := make([]byte, 19)
	b[0] = ipv6
	b[16] = 1
	b[17] = 0
	b[18] = 80
	buf := bytes.NewBuffer(b)
	addr, err := destinationAddr(buf)
	if err != nil {
		t.Errorf("Expected nil, got [%v]", err)
	}
	if addr.String() != "127.0.0.1:80" && addr.String() != "[::1]:80" {
		t.Errorf("Expected addr 127.0.0.1:80 or [::1]:80, got %s", addr.String())
	}
}

func TestDestinationAddrDomain(t *testing.T) {
	b := []byte{domainName, 9, 'l', 'o', 'c', 'a', 'l', 'h', 'o', 's', 't', 0, 80}
	buf := bytes.NewBuffer(b)
	addr, err := destinationAddr(buf)
	if err != nil {
		t.Errorf("Expected nil, got [%v]", err)
	}
	if addr.String() != "127.0.0.1:80" && addr.String() != "[::1]:80" {
		t.Errorf("Expected addr 127.0.0.1:80 or [::1]:80, got %s", addr.String())
	}
}

func TestDestinationAddrNotSupportedType(t *testing.T) {
	b := []byte{5}
	buf := bytes.NewBuffer(b)
	_, err := destinationAddr(buf)
	if err == nil {
		t.Errorf("Expected not supported type error, got nil")
	}
	expected := []byte{5, 8}
	if !bytes.Equal(expected, buf.Bytes()) {
		t.Errorf("Excpected %v, got %v", expected, buf.Bytes())
	}
}
