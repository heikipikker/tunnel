package main

import (
	"net"
	"sync"
	"time"
)

type udpSession struct {
	conn   net.Conn
	live   bool
	from   *net.UDPAddr
	die    chan bool
	readch chan []byte
	clean  func()
}

func newUDPSession(conn net.Conn) *udpSession {
	return &udpSession{
		conn: conn,
		live: true,
		from: conn.RemoteAddr().(*net.UDPAddr),
		die:  make(chan bool),
	}
}

func (sess *udpSession) Close() {
	select {
	case <-sess.die:
	default:
		sess.conn.Close()
		close(sess.die)
		if sess.clean != nil {
			sess.clean()
		}
	}
}

func sessionsCleaner(sessions map[string]*udpSession, lock *sync.Mutex, die chan bool, d time.Duration) {
	ticker := time.NewTicker(d)
	for {
		select {
		case <-die:
			return
		case <-ticker.C:
			var closeSessions []*udpSession
			lock.Lock()
			for k, v := range sessions {
				if v.live {
					v.live = false
				} else {
					delete(sessions, k)
					closeSessions = append(closeSessions, v)
				}
			}
			lock.Unlock()
			for _, v := range closeSessions {
				v.Close()
			}
		}
	}
}

func RunUDPServer(conn net.PacketConn, create func(*SubConn) (net.Conn, net.Conn, error), c *config) {
	defer conn.Close()
	die := make(chan bool)
	defer close(die)
	buf := make([]byte, c.Mtu)
	sessions := make(map[string]*udpSession)
	var lock sync.Mutex

	go sessionsCleaner(sessions, &lock, die, time.Second*time.Duration(c.Expires))

	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			return
		}
		addrstr := addr.String()
		lock.Lock()
		sess, ok := sessions[addrstr]
		lock.Unlock()
		if !ok {
			if create == nil {
				continue
			}
			sconn := &SubConn{
				die:        die,
				sdie:       make(chan bool),
				readch:     make(chan []byte),
				addr:       addr,
				PacketConn: conn,
			}
			var conn1 net.Conn
			var conn2 net.Conn
			conn1, conn2, err = create(sconn)
			if err != nil {
				continue
			}
			sess = &udpSession{
				conn:   sconn,
				live:   true,
				from:   addr.(*net.UDPAddr),
				die:    sconn.sdie,
				clean:  nil,
				readch: sconn.readch,
			}
			lock.Lock()
			sessions[addrstr] = sess
			lock.Unlock()
			go Pipe(conn1, conn2, c.Mtu)
		}
		select {
		case <-sess.die:
		case sess.readch <- buf[:n]:
		}
	}
}

func newUDPListener(address string) (conn *net.UDPConn, err error) {
	laddr, err := net.ResolveUDPAddr("udp", address)
	if err == nil {
		conn, err = net.ListenUDP("udp", laddr)
	}
	return
}

func Pipe(c1, c2 net.Conn, mtu int) {
	c1die := make(chan bool)
	c2die := make(chan bool)
	f := func(dst, src net.Conn, die chan bool, buf []byte) {
		defer close(die)
		var n int
		var err error
		for err == nil {
			n, err = src.Read(buf)
			if n > 0 || err == nil {
				_, err = dst.Write(buf[:n])
			}
		}
	}
	go f(c1, c2, c1die, make([]byte, mtu))
	go f(c2, c1, c2die, make([]byte, mtu))
	select {
	case <-c1die:
	case <-c2die:
	}
}
