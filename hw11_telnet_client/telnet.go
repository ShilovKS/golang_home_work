package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

// Connect устанавливает TCP-соединение с заданным адресом.
func (t *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}
	t.conn = conn
	return nil
}

// Close закрывает TCP-соединение.
func (t *telnetClient) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}

// Send копирует данные из t.in в соединение.
func (t *telnetClient) Send() error {
	// Функция io.Copy блокируется, пока не завершится чтение.
	_, err := io.Copy(t.conn, t.in)
	return err
}

// Receive копирует данные из соединения в t.out.
func (t *telnetClient) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	return err
}
