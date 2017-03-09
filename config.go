package main

import (
	"encoding/json"
	"io/ioutil"
	ss "github.com/ccsexyz/shadowsocks-go/shadowsocks"
)

type config struct {
	Type string `json:"type"`
	Localaddr  string `json:"localaddr"`
	Remoteaddr string `json:"remoteaddr"`
	NoHTTP     bool   `json:"nohttp"`
	Host       string `json:"host"`
	IgnRST     bool   `json:"ignrst"`
	Expires    int    `json:"expires"`
	Method     string `json:"method"`
	Password   string `json:"method"`
	Ivlen int
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
}
