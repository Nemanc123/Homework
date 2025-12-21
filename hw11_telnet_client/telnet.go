package main

import (
	"bufio"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type TelnetClientImpl struct {
	address string
	conn    net.Conn
	in      *bufio.Reader
	out     *bufio.Writer
	timeout time.Duration
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientImpl{
		address: address,
		timeout: timeout,
	}
}
func (t *TelnetClientImpl) Connect() (err error) {
	if t.conn != nil {
		return
	}
	con, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return
	}
	t.conn = con
	t.in = bufio.NewReader(t.conn)
	t.out = bufio.NewWriter(t.conn)
	return
}
func (t *TelnetClientImpl) Send() (err error) {
	if t.conn == nil {
		return
	}
	_, err = io.Copy(t.out, os.Stdin)
	return
}
func (t *TelnetClientImpl) Receive() (err error) {
	if t.conn == nil {
		return
	}
	_, err = io.Copy(os.Stdout, t.in)
	return
}
func (t *TelnetClientImpl) Close() (err error) {
	if t.conn == nil {
		return
	}
	return t.conn.Close()
}
