package socks5

import (
	"net"
	"time"
)

type connection struct {
	conn     net.Conn
	deadline time.Duration
}

func (c *connection) Read(b []byte) (int, error) {
	if c.deadline != 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.deadline))
	}
	return c.conn.Read(b)
}

func (c *connection) Write(b []byte) (int, error) {
	if c.deadline != 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.deadline))
	}
	return c.conn.Write(b)
}
