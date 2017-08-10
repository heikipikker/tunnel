package main

import (
	"log"
	"net"

	"github.com/ccsexyz/rawcon"
)

func RunRemoteServer(c *config) {
	raw := rawcon.Raw{
		NoHTTP: c.NoHTTP,
		Host:   c.Host,
		DSCP:   0,
		IgnRST: true,
	}
	conn, err := raw.ListenRAW(c.Localaddr)
	if err != nil {
		log.Fatal(err)
	}
	create := func(sconn *SubConn) (conn net.Conn, rconn net.Conn, err error) {
		conn = &Conn{
			Conn:   sconn,
			config: c,
		}
		if c.DataShard != 0 && c.ParityShard != 0 {
			conn = &FecConn{
				Conn:       conn,
				config:     c,
				fecEncoder: newFECEncoder(c.DataShard, c.ParityShard, 0),
				fecDecoder: newFECDecoder(3*(c.DataShard+c.ParityShard), c.DataShard, c.ParityShard),
			}
		}
		rconn, err = net.Dial("udp", c.Remoteaddr)
		return
	}
	RunUDPServer(conn, create, c)
}
