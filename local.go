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
	create := func(sconn *SubConn) (conn net.Conn, rconn net.Conn, err error) {
		conn = sconn
		rconn, err = raw.DialRAW(c.Remoteaddr)
		if err != nil {
			return
		}
		rconn = &Conn{
			Conn:   rconn,
			config: c,
		}
		if c.DataShard != 0 && c.ParityShard != 0 {
			rconn = &FecConn{
				Conn:       rconn,
				config:     c,
				fecEncoder: newFECEncoder(c.DataShard, c.ParityShard, 0),
				fecDecoder: newFECDecoder(3*(c.DataShard+c.ParityShard), c.DataShard, c.ParityShard),
			}
		}
		return
	}
	RunUDPServer(conn, create, c)
}
