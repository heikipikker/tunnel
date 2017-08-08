package main

import (
	"log"
	"net"

	"github.com/ccsexyz/rawcon"
)

func RunLocalServer(c *config) {
	raw := rawcon.Raw{
		NoHTTP: c.NoHTTP,
		Host:   c.Host,
		DSCP:   0,
		IgnRST: true,
	}
	conn, err := newUDPListener(c.Localaddr)
	if err != nil {
		log.Fatal(err)
	}
	handle := func(sess *udpSession, b []byte) {
		if len(c.Method) == 0 {
			sess.conn.Write(b)
			return
		}
		buf := make([]byte, 2048)
		n := encrypt(c, b, buf)
		sess.conn.Write(buf[:n])
	}
	create := func(b []byte, from net.Addr) (rconn net.Conn, clean func(), err error) {
		rconn, err = raw.DialRAW(c.Remoteaddr)
		if err == nil {
			if len(c.Method) == 0 {
				rconn.Write(b)
			} else {
				buf := make([]byte, 2048)
				n := encrypt(c, b, buf)
				rconn.Write(buf[:n])
			}
		}
		return
	}
	RunUDPServer(conn, nil, handle, create, c)
}
