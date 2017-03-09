package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	ss "github.com/ccsexyz/shadowsocks-go/shadowsocks"
)

const (
	defaultExpires  = 30
	defaultPassword = "123456"
)

type config struct {
	Type       string `json:"type"`
	Localaddr  string `json:"localaddr"`
	Remoteaddr string `json:"remoteaddr"`
	NoHTTP     bool   `json:"nohttp"`
	Host       string `json:"host"`
	IgnRST     bool   `json:"ignrst"`
	Expires    int    `json:"expires"`
	Method     string `json:"method"`
	Password   string `json:"password"`
	Ivlen      int
}

func readConfig(path string) (configs []*config, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &configs)
	if err != nil {
		var c config
		err = json.Unmarshal(bytes, &c)
		if err == nil {
			configs = append(configs, &c)
		}
	}
	for _, c := range configs {
		checkConfig(c)
	}
	return
}

func checkConfig(c *config) {
	c.Ivlen = ss.GetIvLen(c.Method)
	if c.Expires == 0 {
		c.Expires = defaultExpires
	}
	if len(c.Method) != 0 && len(c.Password) == 0 {
		c.Password = defaultPassword
	}
	if len(c.Localaddr) == 0 {
		log.Fatal("no localaddr")
	}
	if len(c.Remoteaddr) == 0 {
		log.Fatal("no remoteaddr")
	}
}
