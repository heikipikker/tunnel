package main

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	LocalAddr  string `json:"localaddr"`
	RemoteAddr string `json:"remoteaddr"`
	NoHTTP     bool   `json:"nohttp"`
	Host       string `json:"host"`
	IgnRST     bool   `json:"ignrst"`
}

func readCofnig(path string) (configs []config, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &configs)
	return
}
