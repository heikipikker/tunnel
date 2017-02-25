package main

import (
	"net"
	"sync"
	"time"

	"github.com/ccsexyz/rawcon"
)

type client struct {
	c   config
	die chan bool
	r   rawcon.Raw
}

func newClient(c config) *client {
	return &client{
		c:   c,
		die: make(chan bool),
		r: rawcon.Raw{
			NoHTTP: c.NoHTTP,
			IgnRST: c.IgnRST,
			Host:   c.Host,
			DSCP:   0,
		},
	}
}

func (c *client) createListener() *net.UDPConn {
	laddr, err := net.ResolveUDPAddr("udp", c.c.LocalAddr)
	fatalErr(err)
	conn, err := net.ListenUDP("udp", laddr)
	fatalErr(err)
	return conn
}

func (c *client) createRawConn() (conn net.Conn, err error) {
	conn, err = c.r.DialRAW(c.c.RemoteAddr)
	return
}

func (c *client) close() {
	close(c.die)
}

func (c *client) run() {
	defer c.close()
	for {
		c.runOnce()
	}
}

func (c *client) runOnce() {
	die := make(chan bool)
	buf := make([]byte, 65536)
	var mutex sync.Mutex
	sessions := make(map[string]net.Conn)
	newsessions := make(map[string]*time.Time)
	listener := c.createListener()
	defer close(die)
	go func() {
		defer listener.Close()
		defer func() {
			mutex.Lock()
			defer mutex.Unlock()
			for _, v := range sessions {
				v.Close()
			}
		}()
		select {
		case <-die:
		case <-c.die:
		}
	}()
	writefunc := func(conn net.Conn, addrstr string) {
		defer conn.Close()
		addr, err := net.ResolveUDPAddr("udp", addrstr)
		if err != nil {
			return
		}
		rbuf := make([]byte, 65536)
		for {
			n, err := conn.Read(rbuf)
			if err != nil {
				select {
				case <-die:
				case <-c.die:
				default:
					mutex.Lock()
					now := time.Now()
					newsessions[addrstr] = &now
					mutex.Unlock()
				}
				return
			}
			_, err = listener.WriteTo(rbuf[:n], addr)
			if err != nil {
				listener.Close()
				return
			}
		}
	}
	for {
		n, addr, err := listener.ReadFromUDP(buf)
		if err != nil {
			return
		}
		addrstr := addr.String()
		mutex.Lock()
		conn, ok := sessions[addrstr]
		t, ok2 := newsessions[addrstr]
		mutex.Unlock()
		if ok {
			_, err = conn.Write(buf[:n])
			if err != nil {
				conn.Close()
			}
			continue
		}
		now := time.Now()
		if ok2 {
			// drop this packet
			if time.Now().After(t.Add(time.Minute)) {
				mutex.Lock()
				delete(newsessions, addrstr)
				mutex.Unlock()
			}
			continue
		}
		mutex.Lock()
		newsessions[addrstr] = &now
		mutex.Unlock()
		var wbuf []byte
		if n > 0 {
			wbuf = make([]byte, n)
			copy(wbuf, buf)
		}
		go func(str string, buf []byte) {
			// createRawConn is a block function
			// so start a new goroutine
			conn, err := c.createRawConn()
			if err != nil {
				return
			}
			select {
			case <-c.die:
				conn.Close()
				return
			default:
				_, err = conn.Write(wbuf)
				if err != nil {
					conn.Close()
					return
				}
				mutex.Lock()
				delete(newsessions, addrstr)
				sessions[addrstr] = conn
				mutex.Unlock()
				go writefunc(conn, addrstr)
			}
		}(addrstr, wbuf)
	}
}
