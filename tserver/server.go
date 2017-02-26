package main

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/ccsexyz/rawcon"
)

type server struct {
	c   config
	die chan bool
	r   rawcon.Raw
}

func newServer(c config) *server {
	return &server{
		c:   c,
		die: make(chan bool),
		r: rawcon.Raw{
			NoHTTP: c.NoHTTP,
			IgnRST: c.IgnRST,
			DSCP:   0,
		},
	}
}

func (s *server) createRawListener() (listener net.PacketConn) {
	listener, err := s.r.ListenRAW(s.c.LocalAddr)
	fatalErr(err)
	return
}

func (s *server) createUDPConn() (conn net.Conn, err error) {
	raddr, err := net.ResolveUDPAddr("udp", s.c.TargetAddr)
	fatalErr(err)
	conn, err = net.DialUDP("udp", nil, raddr)
	return
}

func (s *server) close() {
	close(s.die)
}

func (s *server) run() {
	defer s.close()
	for {
		s.runOnce()
	}
}

func (s *server) runOnce() {
	die := make(chan bool)
	buf := make([]byte, 65536)
	var mutex sync.Mutex
	sessions := make(map[string]net.Conn)
	expires := make(map[string]bool)
	listener := s.createRawListener()
	defer close(die)
	go func() {
		defer log.Println("listener was closed")
		defer listener.Close()
		defer func() {
			mutex.Lock()
			defer mutex.Unlock()
			for _, v := range sessions {
				v.Close()
			}
		}()
		ticker := time.NewTicker(time.Second * time.Duration(s.c.Expires))
		for {
			select {
			case <-die:
				return
			case <-s.die:
				return
			case <-ticker.C:
				mutex.Lock()
				for k, v := range expires {
					if v {
						expires[k] = !v
					} else {
						delete(expires, k)
						conn, ok := sessions[k]
						if ok {
							delete(sessions, k)
							conn.Close()
						}
					}
				}
				mutex.Unlock()
			}
		}
	}()
	readfunc := func(conn net.Conn, addrstr string) {
		defer conn.Close()
		defer func() {
			select {
			case <-die:
			case <-s.die:
			default:
				mutex.Lock()
				delete(sessions, addrstr)
				delete(expires, addrstr)
				mutex.Unlock()
			}
		}()
		addr, err := net.ResolveUDPAddr("udp", addrstr)
		if err != nil {
			return
		}
		buf := make([]byte, 65536)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				return
			}
			_, err = listener.WriteTo(buf[:n], addr)
			if err != nil {
				return
			}
		}
	}
	for {
		n, addr, err := listener.ReadFrom(buf)
		if err != nil {
			return
		}
		addrstr := addr.String()
		mutex.Lock()
		conn, ok := sessions[addrstr]
		expires[addrstr] = true
		mutex.Unlock()
		if ok {
			_, err = conn.Write(buf[:n])
			if err != nil {
				conn.Close()
			}
			continue
		}
		conn, err = s.createUDPConn()
		fatalErr(err)
		_, err = conn.Write(buf[:n])
		if err != nil {
			conn.Close()
			continue
		}
		mutex.Lock()
		sessions[addrstr] = conn
		mutex.Unlock()
		go readfunc(conn, addrstr)
	}
}
