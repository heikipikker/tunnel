package main

import (
	"log"
	"net"
	"sync"
	"time"

	ss "github.com/ccsexyz/shadowsocks-go/shadowsocks"
)

type udpSession struct {
	conn  net.Conn
	live  bool
	from  *net.UDPAddr
	die   chan bool
	clean func()
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

func RunUDPServer(conn net.PacketConn, check func([]byte) bool, handle func(*udpSession, []byte),
	create func([]byte, net.Addr) (net.Conn, func(), error), c *config) {
	defer conn.Close()
	die := make(chan bool)
	defer close(die)
	buf := make([]byte, 2048)
	dbuf := make([]byte, 2048)
	sessions := make(map[string]*udpSession)
	var lock sync.Mutex

	go sessionsCleaner(sessions, &lock, die, time.Second*time.Duration(c.Expires))

	for {
		rbuf := buf
		n, addr, err := conn.ReadFrom(rbuf)
		if err != nil {
			return
		}
		if len(c.Method) != 0 && c.Type == "server" {
			if n < c.Ivlen {
				continue
			}
			n = decrypt(c, rbuf[:n], dbuf)
			rbuf = dbuf
		}
		if check != nil && !check(rbuf[:n]) {
			continue
		}
		addrstr := addr.String()
		lock.Lock()
		sess, ok := sessions[addrstr]
		lock.Unlock()
		if ok {
			sess.live = true
			if handle != nil {
				handle(sess, rbuf[:n])
			}
		} else {
			if create != nil {
				rconn, clean, err := create(rbuf[:n], addr)
				if err != nil {
					log.Println(err)
					continue
				}
				if rconn == nil {
					continue
				}
				sess = &udpSession{conn: rconn, live: true, from: addr.(*net.UDPAddr), die: make(chan bool), clean: clean}
				lock.Lock()
				sessions[addrstr] = sess
				lock.Unlock()
				go func(sess *udpSession) {
					defer sess.Close()
					buf := make([]byte, 2048)
					ebuf := make([]byte, 2048)
					for {
						wbuf := buf[:]
						n, err := sess.conn.Read(wbuf)
						if err != nil {
							return
						}
						if len(c.Method) != 0 {
							if c.Type == "server" {
								n = encrypt(c, wbuf[:n], ebuf)
							} else {
								n = decrypt(c, wbuf[:n], ebuf)
							}
							wbuf = ebuf
						}
						_, err = conn.WriteTo(wbuf[:n], sess.from)
					}
				}(sess)
			}
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

func encrypt(c *config, buf, ebuf []byte) (n int) {
	enc, err := ss.NewEncrypter(c.Method, c.Password)
	if err != nil {
		log.Fatal(err)
	}
	n = copy(ebuf, enc.GetIV())
	enc.Encrypt(ebuf[n:], buf)
	n += len(buf)
	return
}

func decrypt(c *config, buf, dbuf []byte) (n int) {
	ivlen := c.Ivlen
	if len(buf) < ivlen {
		return -1
	}
	dec, err := ss.NewDecrypter(c.Method, c.Password, buf[:ivlen])
	if err != nil {
		return -1
	}
	dec.Decrypt(dbuf, buf[ivlen:])
	n = len(buf) - ivlen
	return
}
