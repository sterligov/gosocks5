package socks5

import (
	"bytes"
	"testing"
)

func TestNoAuth(t *testing.T) {
	b := []byte{2,0,1}
	buf := bytes.NewBuffer(b)
	err := authorize(nil, buf)
	if err != nil {
		t.Errorf("Excpected nil error, got [%s]", err.Error())
	}
	expected := []byte{5, 0}
	if !bytes.Equal(expected, buf.Bytes()) {
		t.Errorf("Excpected %v, got %v", expected, buf.Bytes())
	}
}

func TestNoAcceptableAuthMethods(t *testing.T) {
	b := []byte{1,3}
	buf := bytes.NewBuffer(b)
	err := authorize(nil, buf)
	if err == nil {
		t.Errorf("Excpected not acceptable error, got nil")
	}
	expected := []byte{5, 255}
	if !bytes.Equal(expected, buf.Bytes()) {
		t.Errorf("Excpected %v, got %v", expected, buf.Bytes())
	}
}

func TestBadLogin(t *testing.T) {
	b := []byte{2,0,2,5,4,'u','s','e','r',3,'b','a','d'}
	buf := bytes.NewBuffer(b)
	err := authorize(&RoughAuthorizer{"user", "pass"}, buf)
	if err == nil {
		t.Errorf("Excpected bad login details error, got nil")
	}
	expected := []byte{5,2,1,1}
	if !bytes.Equal(expected, buf.Bytes()) {
		t.Errorf("Excpected %v, got %v", expected, buf.Bytes())
	}
}

func TestCorrectLogin(t *testing.T) {
	b := []byte{2,0,2,5,4,'u','s','e','r',4,'p','a','s','s'}
	buf := bytes.NewBuffer(b)
	err := authorize(&RoughAuthorizer{"user", "pass"}, buf)
	if err != nil {
		t.Errorf("Excpected nil, got [%v]", err)
	}
	expected := []byte{5,2,1,0}
	if !bytes.Equal(expected, buf.Bytes()) {
		t.Errorf("Excpected %v, got %v", expected, buf.Bytes())
	}
}