package tcp

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

func NewConnect(server *Server, conn *net.TCPConn, addr *Addr) *Connect {
	c := new(Connect)
	c.server = server
	c.conn = conn
	c.addr = addr
	c.recvCh = make(chan *Message, 10)
	c.sendCh = make(chan *Message, 10)
	c.running = true
	return c
}

type Connect struct {
	name      string
	server    *Server
	addr      *Addr
	conn      *net.TCPConn
	recvCh    chan *Message
	sendCh    chan *Message
	running   bool
	wg        sync.WaitGroup
	closeLock sync.Mutex
}

func (c *Connect) worker() {
	go c.recv()
	go c.handle()
	go c.send()
}

func (c *Connect) recv() {
	c.wg.Add(1)
	defer c.wg.Done()

	for c.running {
		time.Sleep(time.Second)
		msg := Message{}

		err := read(c.conn, &msg)
		if err != nil {
			switch err {
			case io.EOF /*errRecvEOF, errRemoteForceDisconnect*/ :
				fmt.Printf("this client connect is close: %s.", err.Error())
				c.server.closeConnect(c.addr)
				c.Close()
				return
			default:
				fmt.Printf("recv msg err: %s.", err.Error())
			}

			continue
		}

		c.recvCh <- &msg
	}
}

func (c *Connect) handle() {
	c.wg.Add(1)
	defer c.wg.Done()

	for c.running {
		select {
		case msg := <-c.recvCh:
			go func() {
				c.server.callOnRecv(c.addr, msg)
			}()
		}
	}
}

func (c *Connect) send() {
	c.wg.Add(1)
	defer c.wg.Done()

	for c.running {
		select {
		case msg := <-c.sendCh:
			data, err := pack(msg)
			if err != nil {
				fmt.Printf("pack data address(%s) error:%s.", c.addr.GetAddress(), err.Error())
				continue
			}

			_, err = c.conn.Write(data.Bytes())
			if err != nil {
				fmt.Printf("send data to client address(%s) error:%s.", c.addr.GetAddress(), err.Error())
				c.server.closeConnect(c.addr)
				c.Close()
				return
			}
		}
	}
}

func (c *Connect) Close() {
	c.closeLock.Lock()
	defer c.closeLock.Unlock()
	if c.running {
		c.running = false
		c.conn.Close()
		c.wg.Wait()
		close(c.recvCh)
		close(c.sendCh)
	}
}
