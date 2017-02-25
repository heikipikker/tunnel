package main

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	LocalAddr  string `json:"localaddr"`
	TargetAddr string `json:"targetaddr"`
	NoHTTP     bool   `json:"nohttp"`
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
