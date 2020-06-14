package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var (
	ErrEOF          = errors.New("...EOF")
	ErrClosedByPeer = errors.New("...Connection was closed by peer")
)

type TelnetClient interface {
	Connect() error
	Send() error
	Receive() error
	Close() error
}

type BasicTelnetClient struct {
	address  string
	timeout  time.Duration
	conn     net.Conn
	connScan *bufio.Scanner
	inScan   *bufio.Scanner
	out      io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient { //nolint
	c := &BasicTelnetClient{
		address: address,
		timeout: timeout,
		inScan:  bufio.NewScanner(in),
		out:     out,
	}
	return c
}

func (c *BasicTelnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	c.connScan = bufio.NewScanner(c.conn)
	return nil
}

func (c *BasicTelnetClient) Send() error {
	if !c.inScan.Scan() {
		if c.inScan.Err() == nil {
			// здесь проще вывести сообщение об этой ошибке - больше контекста
			fmt.Fprintf(os.Stderr, "%s\n", ErrEOF)
			return ErrEOF
		}
		return c.inScan.Err()
	}
	str := c.inScan.Text()
	_, err := c.conn.Write([]byte(fmt.Sprintf("%s\n", str)))
	return err
}

func (c *BasicTelnetClient) Receive() error {
	if !c.connScan.Scan() {
		if c.connScan.Err() == nil {
			fmt.Fprintf(os.Stderr, "%s\n", ErrClosedByPeer)
			return ErrClosedByPeer
		}
		return c.connScan.Err()
	}
	str := c.connScan.Text()
	_, err := c.out.Write([]byte(fmt.Sprintf("%s\n", str)))
	return err
}

func (c *BasicTelnetClient) Close() error {
	return c.conn.Close()
}
