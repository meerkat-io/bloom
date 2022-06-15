package tcp

import (
	"net"
)

type Server interface {
	Accept(*Connection)
}

type Listener struct {
	socket *net.TCPListener
	server Server
}

func Listen(addr string, server Server) (*Listener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	socket, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, err
	}
	l := &Listener{
		socket: socket,
		server: server,
	}
	go l.listen()
	return l, nil
}

func (l *Listener) Close() {
	l.socket.Close()
}

func (l *Listener) listen() {
	for {
		if socket, err := l.socket.AcceptTCP(); err == nil {
			conn := &Connection{
				socket: socket,
			}
			go l.server.Accept(conn)
		} else {
			return
		}
	}
}
