package socks5

import (
	"fmt"
	"io"
)

const (
	noAuth                   = 0x00
	usernameAuth             = 0x02
	noAcceptableAuthMethods  = 0xff
	authVer                  = 0x01
	authOK                   = 0x00
	authFail                 = 0x01
)

type Authorizer interface {
	Authorize(string, string) bool
}

type RoughAuthorizer struct {
	Username string
	Password string
}

func (r *RoughAuthorizer) Authorize(username, password string) bool {
	return r.Username == username && r.Password == password
}

func authorize(authorizer Authorizer, rw io.ReadWriter) error {
	methods, err := authMethods(rw)
	if err != nil {
		return err
	}

	if _, ok := methods[noAuth]; ok && authorizer == nil {
		rw.Write([]byte{socksV5, noAuth})
		return nil
	} else if _, ok := methods[usernameAuth]; ok && authorizer != nil {
		rw.Write([]byte{socksV5, usernameAuth})
	} else {
		rw.Write([]byte{socksV5, noAcceptableAuthMethods})
		return fmt.Errorf("no acceptable auth methods")
	}

	username, password, err := loginDetails(rw)
	if err != nil {
		return err
	}

	if !authorizer.Authorize(username, password) {
		rw.Write([]byte{authVer, authFail})
		return fmt.Errorf("bad login details")
	}
	rw.Write([]byte{authVer, authOK})

	return nil
}

func loginDetails(r io.Reader) (string, string, error) {
	buf := make([]byte, 2)
	_, err := r.Read(buf)
	if err != nil {
		return "", "", err
	}

	username := make([]byte, buf[1])
	_, err = r.Read(username)
	if err != nil {
		return "", "", err
	}

	buf = make([]byte, 1)
	_, err = r.Read(buf)
	if err != nil {
		return "", "", err
	}

	password := make([]byte, buf[0])
	_, err = r.Read(password)
	if err != nil {
		return "", "", err
	}

	return string(username), string(password), nil
}

func authMethods(r io.Reader) (map[byte]struct{}, error) {
	buf := make([]byte, 1)
	_, err := r.Read(buf)
	if err != nil {
		return nil, err
	}

	methods := make([]byte, buf[0])
	_, err = r.Read(methods)
	if err != nil {
		return nil, err
	}

	m := make(map[byte]struct{}, buf[0])
	for _, v := range methods {
		m[v] = struct{}{}
	}

	return m, nil
}
