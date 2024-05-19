package server

import (
	"fmt"
	"goTask/tcp"
	"net"
)

func main() {
	server := tcp.NewServer()

	// on connect event
	server.OnConnect(func(conn *net.TCPConn, addr *tcp.Addr) {
		fmt.Println(fmt.Sprintf("one client connect, remote address=%s.", conn.RemoteAddr().String()))
	})

	// on receive data event
	server.OnRecv(func(addr *tcp.Addr, req *tcp.Message) {
		fmt.Println(fmt.Sprintf("req.Type=%d, req.Data=%s.", req.Type, string(req.Data)))

		_ = server.Send(*addr, &tcp.Message{
			Type: 1,
			Data: []byte(fmt.Sprintf("hello: %s.", addr.GetAddress())),
		})
	})

	// on disconnect event
	server.OnDisconnect(func(addr *tcp.Addr) {

	})

	server.Run(":8080")
}
