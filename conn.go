package main

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/ccsexyz/utils"
)

type Conn struct {
	net.Conn
	*config
}

func (c *Conn) Read(b []byte) (n int, err error) {
	b2 := make([]byte, c.Mtu)
	for {
		n, err = c.Conn.Read(b2)
		if err != nil {
			return
		}
		if n <= c.Ivlen || n > c.Mtu {
			continue
		}
		var dec utils.Decrypter
		dec, err = utils.NewDecrypter(c.Method, c.Password, b2[:c.Ivlen])
		if err != nil {
			return
		}
		dec.Decrypt(b, b2[c.Ivlen:n])
		n -= c.Ivlen
		return
	}
}

func (c *Conn) Write(b []byte) (n int, err error) {
	defer func() {
		// log.Println(n, err)
	}()
	n2 := len(b) + c.Ivlen
	if n2 > c.Mtu {
		err = fmt.Errorf("buffer is too large")
		return
	}
	enc, err := utils.NewEncrypter(c.Method, c.Password)
	if err != nil {
		return
	}
	b2 := make([]byte, c.Mtu)
	copy(b2, enc.GetIV())
	enc.Encrypt(b2[c.Ivlen:], b)
	_, err = c.Conn.Write(b2[:n2])
	if err == nil {
		n = len(b)
	}
	return
}

// FecConn implements FEC decoder and encoder
type FecConn struct {
	net.Conn
	*config
	fecDecoder *fecDecoder
	fecEncoder *fecEncoder
	recovers   [][]byte
}

func (c *FecConn) Read(b []byte) (n int, err error) {
	for n == 0 {
		for len(c.recovers) != 0 {
			r := c.recovers[0]
			c.recovers = c.recovers[1:]
			if len(r) < 2 {
				continue
			}
			sz := int(binary.LittleEndian.Uint16(r))
			if sz < 2 || sz > len(r) {
				continue
			}
			n = copy(b, r[2:sz])
			return
		}
		buf := make([]byte, c.Mtu)
		var num int
		num, err = c.Conn.Read(buf)
		if err != nil {
			return
		}
		f := c.fecDecoder.decodeBytes(buf)
		if f.flag == typeData {
			n = copy(b, buf[fecHeaderSizePlus2:num])
		}
		if f.flag == typeData || f.flag == typeFEC {
			c.recovers = c.fecDecoder.decode(f)
		}
	}
	return
}

func (c *FecConn) Write(b []byte) (n int, err error) {
	ext := b[:fecHeaderSizePlus2+len(b)]
	copy(ext[fecHeaderSizePlus2:], b)
	ecc := c.fecEncoder.encode(ext)

	n, err = c.Conn.Write(ext)
	if err != nil {
		return
	}

	for _, e := range ecc {
		_, err = c.Conn.Write(e)
		if err != nil {
			return
		}
	}

	return
}
