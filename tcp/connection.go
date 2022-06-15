package tcp

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	packetHeaderSize = 4
	maxPacketSize    = 1024 * 1024
)

// Connection data structure
type Connection struct {
	socket        *net.TCPConn
	maxPacketSize int
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
		socket:        conn,
		maxPacketSize: maxPacketSize,
	}, nil
}

// Set max packet size (default: 1024 * 1024)
func (c *Connection) SetMaxPacketSize(size int) {
	c.maxPacketSize = size
}

// Receive data from socket
func (c *Connection) Receive() ([]byte, error) {
	var data []byte
	length, err := c.readPacketSize()
	if err == nil {
		readed := 0
		bytes := 0
		data = make([]byte, length)
		for readed < length && err == nil {
			bytes, err = c.socket.Read(data[readed:])
			readed += bytes
		}
	}
	return data, err
}

// Send data to socket
func (c *Connection) Send(data []byte) error {
	length := len(data)
	err := c.writePacketSize(length)
	if err == nil {
		writed := 0
		bytes := 0
		for writed < length && err == nil {
			bytes, err = c.socket.Write(data[writed:])
			writed += bytes
		}
	}
	return err
}

// Close connection
func (c *Connection) Close() {
	c.socket.Close()
}

func (c *Connection) readPacketSize() (int, error) {
	header := make([]byte, packetHeaderSize)
	size, err := c.socket.Read(header)
	if err != nil || size != packetHeaderSize {
		return 0, fmt.Errorf("read socket error from %s", c.socket.RemoteAddr())
	}
	length := int(binary.LittleEndian.Uint32(header))
	if length == 0 {
		return 0, fmt.Errorf("empty packet from %s", c.socket.RemoteAddr())
	}
	if length > c.maxPacketSize {
		return 0, fmt.Errorf("data overflow from %s", c.socket.RemoteAddr())
	}
	return length, nil
}

func (c *Connection) writePacketSize(length int) error {
	if length > c.maxPacketSize {
		return fmt.Errorf("data overflow write to %s", c.socket.RemoteAddr())
	}
	header := make([]byte, packetHeaderSize)
	binary.LittleEndian.PutUint32(header, uint32(length))
	bytesWrite, err := c.socket.Write(header)
	if err != nil || bytesWrite != packetHeaderSize {
		return fmt.Errorf("write data error to %s", c.socket.RemoteAddr())
	}
	return nil
}
