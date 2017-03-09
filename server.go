package main

import (
	"log"
	"net"
	"github.com/ccsexyz/rawcon"
)

func RunRemoteServer(c *config) {
	raw := rawcon.Raw{
		NoHTTP: c.NoHTTP,
		Host:  c.Host,
		DSCP: 0,
		IgnRST: c.IgnRST,
	}
	conn, err := raw.ListenRAW(c.Localaddr)
	if err != nil {
		log.Fatal(err)
	}
	handle := func(sess *udpSession, b []byte) {
		sess.conn.Write(b)
	}
	create := func(b []byte, from net.Addr) (rconn net.Conn, clean func(), err error) {
		rconn, err = net.Dial("udp", c.Remoteaddr)
		if err == nil {
			rconn.Write(b)
		}
		return
	}
	RunUDPServer(conn, nil, handle, create, c)
}
