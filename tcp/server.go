package tcp

import (
	"fmt"
	"net"
	"sync"
)

type ServerRecvHandle = func(addr *Addr, req *Message)
type ConnectHandle = func(conn *net.TCPConn, addr *Addr)
type DisconnectHandle = func(addr *Addr)

func NewServer() *Server {
	s := new(Server)

	s.connMgr = make(map[Addr]*Connect)
	s.running = true
	return s
}

type Server struct {
	listener     *net.TCPListener  // the tcp server Listener
	connMgr      map[Addr]*Connect // client connect map
	onRecv       ServerRecvHandle  // receive event callback function
	onConnect    ConnectHandle     // connect event callback function
	onDisconnect DisconnectHandle  // disconnect event callback function
	lock         sync.RWMutex
	running      bool
}

func (s *Server) OnRecv(handle ServerRecvHandle) {
	s.onRecv = handle
}

func (s *Server) callOnRecv(addr *Addr, req *Message) {
	go func() {
		if s.onRecv != nil {
			s.onRecv(addr, req)
		}
	}()
}

func (s *Server) Send(addr Addr, msg *Message) error {
	if connect, ok := s.connMgr[addr]; !ok {
		return fmt.Errorf("the address[%s] of connect is not exist.", addr.GetAddress())
	} else {
		go func() {
			connect.sendCh <- msg
		}()
	}
	return nil
}

func (s *Server) OnConnect(handle ConnectHandle) {
	s.onConnect = handle
}

func (s *Server) callOnConnect(conn *net.TCPConn, addr *Addr) {
	go func() {
		if s.onConnect != nil {
			s.onConnect(conn, addr)
		}
	}()
}

func (s *Server) OnDisconnect(handle DisconnectHandle) {
	s.onDisconnect = handle
}

func (s *Server) callOnDisconnect(addr *Addr) {
	go func() {
		if s.onDisconnect != nil {
			s.onDisconnect(addr)
		}
	}()
}

func (s *Server) Close() {
	s.lock.Lock()
	defer s.lock.Lock()

	s.running = false

	if s.listener != nil {
		s.listener.Close()
	}

	var wg sync.WaitGroup
	for addr, connect := range s.connMgr {
		go func() {
			wg.Add(1)
			connect.Close()
			delete(s.connMgr, addr)
			wg.Done()
		}()
	}

	wg.Wait()
}

func (s *Server) closeConnect(addr *Addr) {
	s.callOnDisconnect(addr)

	if connect, ok := s.connMgr[*addr]; ok {
		connect.Close()
		delete(s.connMgr, *addr)
	}
}

func (s *Server) Run(addr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		//fmt.Printf("the server listen address err:%s.", err.Error())
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		//fmt.Printf("tcp server listen addr(%s) error:%s.", addr, err.Error())
		return
	}

	s.listener = listener

	for s.running {
		//fmt.Printf("server started waiting for connection")
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			//fmt.Printf("tcp server accept one client error:%s", err.Error())
			continue
		}

		addr := NewAddr(conn.RemoteAddr().String())

		connect := NewConnect(s, conn, addr)

		if s.connMgr == nil {
			s.connMgr = make(map[Addr]*Connect)
		}
		s.connMgr[*addr] = connect
		s.callOnConnect(conn, addr)
		connect.worker()
	}
}
