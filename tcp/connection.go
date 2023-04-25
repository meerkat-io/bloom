package tcp

import (
	"io"
	"net"
)

// Connection data structure
type Connection struct {
	socket *net.TCPConn
}

// Dial an address to create a connection
func Dial(addr string) (*Connection, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	return &Connection{
		socket: conn,
	}, nil
}

// Return io.Reader interface
func (c *Connection) Reader() io.Reader {
	return c.socket
}

// Return io.Writer interface
func (c *Connection) Writer() io.Writer {
	return c.socket
}

// Close connection
func (c *Connection) Close() {
	c.socket.Close()
}
