package main

import (
	"log"
	"net"

	"github.com/ccsexyz/rawcon"
	"github.com/ccsexyz/utils"
)

// RunRemoteServer run the remote server
func RunRemoteServer(c *config) {
	raw := rawcon.Raw{
		Host:   c.Host,
		DSCP:   0,
		IgnRST: true,
		Mixed:  true,
	}
	var conn net.PacketConn
	var err error
	if c.UDP {
		conn, err = utils.NewUDPListener(c.Localaddr)
	} else {
		conn, err = raw.ListenRAW(c.Localaddr)
	}
	if err != nil {
		log.Fatal(err)
	}
	create := func(sconn *utils.SubConn) (conn net.Conn, rconn net.Conn, err error) {
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
	ctx := &utils.UDPServerCtx{Expires: c.Expires, Mtu: c.Mtu}
	ctx.RunUDPServer(conn, create)
}
