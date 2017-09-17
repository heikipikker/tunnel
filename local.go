package main

import (
	"log"
	"net"

	"github.com/ccsexyz/rawcon"
	"github.com/ccsexyz/utils"
)

// RunLocalServer run the local client server
func RunLocalServer(c *config) {
	raw := rawcon.Raw{
		NoHTTP: c.NoHTTP,
		Host:   c.Host,
		DSCP:   0,
		IgnRST: true,
	}
	conn, err := utils.NewUDPListener(c.Localaddr)
	if err != nil {
		log.Fatal(err)
	}
	create := func(sconn *utils.SubConn) (conn net.Conn, rconn net.Conn, err error) {
		conn = sconn
		if c.UDP {
			rconn, err = net.Dial("udp", c.Remoteaddr)
		} else {
			rconn, err = raw.DialRAW(c.Remoteaddr)
		}
		if err != nil {
			log.Println(err)
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
	ctx := &utils.UDPServerCtx{Expires: c.Expires, Mtu: c.Mtu}
	ctx.RunUDPServer(conn, create)
}
