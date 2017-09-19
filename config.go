package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"reflect"

	"github.com/ccsexyz/utils"
)

const (
	defaultExpires  = 30
	defaultPassword = "123456"
)

type config struct {
	Type        string `json:"type"`
	Localaddr   string `json:"localaddr"`
	Remoteaddr  string `json:"remoteaddr"`
	NoHTTP      bool   `json:"nohttp"`
	Host        string `json:"host"`
	Expires     int    `json:"expires"`
	DataShard   int    `json:"datashard"`
	ParityShard int    `json:"parityshard"`
	Method      string `json:"method"`
	Password    string `json:"password"`
	Mtu         int    `json:"mtu"`
	UDP         bool   `json:"udp"`
	Dummy       bool   `json:"dummy"`
	UseMul      bool   `json:"usemul"`
	MulConn     int    `json:"mulconn"`
	Ivlen       int
}

func (c *config) valid() bool {
	if len(c.Localaddr) == 0 {
		return false
	}
	if len(c.Remoteaddr) == 0 && c.Type != "server" {
		return false
	}
	if c.DataShard < 0 || c.ParityShard < 0 || c.Mtu < 0 || c.Expires < 0 {
		return false
	}
	return true
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
	c.Ivlen = utils.GetIvLen(c.Method)
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
	if c.Mtu <= 0 {
		c.Mtu = 1420
	}
}

func (c *config) print() {
	val := reflect.ValueOf(c)
	typ := reflect.Indirect(val).Type()
	nfield := typ.NumField()
	for i := 0; i < nfield; i++ {
		jv := typ.Field(i).Tag.Get("json")
		if len(jv) != 0 {
			log.Println(jv+":", val.Elem().Field(i))
		}
	}
	log.Println("")
}
