package tcp

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

var ()

func init() {

}

type ClientRecvHandle = func(recv *Message)

func NewClient(address string) *Client {
	client := new(Client)
	client.connectServer(address)
	client.running = true
	go client.reconnect()

	return client
}

type Client struct {
	remoteAddr   *Addr            // remote address
	localAddr    *Addr            // local address
	conn         net.Conn         // connect server obj, receive chan, send chan
	onRecv       ClientRecvHandle // receive event callback function
	onDisconnect DisconnectHandle // disconnect event callback function
	running      bool             // is running flag
	connected    bool             // is connect flag
	wg           sync.WaitGroup   // wait event obj, use save exit
	connectLock  sync.Mutex       // connect lock, Avoid concurrent connections
}

func (c *Client) connectServer(address string) error {
	c.connectLock.Lock()
	defer c.connectLock.Unlock()

	if c.connected {
		return nil
	}
	c.remoteAddr = NewAddr(address)

	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return fmt.Errorf("get address info error:%s.", err.Error())
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return fmt.Errorf("connect server error:%s.", err.Error())
	}

	if err = conn.SetKeepAlive(true); err != nil {
		//fmt.Printf("set keep alive err: %s.", err.Error())
	}
	if err = conn.SetKeepAlivePeriod(5 * time.Second); err != nil {
		//fmt.Printf("set keep alive period err: %s.", err.Error())
	}

	c.conn = conn
	c.connected = true
	c.localAddr = NewAddr(conn.LocalAddr().String())

	//fmt.Printf("connect server success, address=%s.", address)
	go c.recv()

	return nil
}

func (c *Client) OnRecv(handle ClientRecvHandle) {
	c.onRecv = handle
}

func (c *Client) callOnRecv(recv *Message) {
	go func() {
		if c.onRecv != nil {
			c.onRecv(recv)
		}
	}()
}

func (c *Client) reconnect() {
	c.wg.Add(1)
	defer c.wg.Done()

	var currentConnectStatus error
	for c.running {
		time.Sleep(time.Second * 3)

		err := c.connectServer(c.remoteAddr.GetAddress())
		if err != nil {
			if currentConnectStatus == nil || err.Error() != currentConnectStatus.Error() {
				currentConnectStatus = err
				fmt.Printf(err.Error())
			}
		}
	}
}

func (c *Client) Send(msg *Message) error {
	if c.conn != nil {
		buf, err := pack(msg)
		if err != nil {
			return fmt.Errorf("pack data error:%s.", err.Error())
		}

		_, err = c.conn.Write(buf.Bytes())
		if err != nil {
			c.connected = false
			return fmt.Errorf("send data error:%s.", err.Error())
		}
	}
	return nil
}

func (c *Client) recv() {
	c.wg.Add(1)
	defer c.wg.Done()

	for c.running && c.connected {
		msg := Message{}
		if err := read(c.conn, &msg); err != nil {
			switch err {
			case io.EOF:
				fmt.Printf("this client connect disconnect: %s.", err.Error())
				c.connected = false
				break
			default:
				fmt.Printf("recv msg err: %s.", err.Error())
			}
			continue
		}

		c.callOnRecv(&msg)
	}
}

func (c *Client) Close() {
	c.running = false
	if c.conn != nil {
		c.conn.Close()
	}
	c.wg.Wait()
}
