package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ccsexyz/utils"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.Lmicroseconds)
	if len(os.Args) != 2 {
		fmt.Println("usage: tunnel configfile")
		return
	}
	configs, err := readConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	for _, c := range configs {
		wg.Add(1)
		go func(c *config) {
			if c.Pprof != "" {
				utils.RunProfileHTTPServer(c.Pprof)
			}
			defer wg.Done()
			switch c.Type {
			case "server":
				RunRemoteServer(c)
			default:
				RunLocalServer(c)
			}
		}(c)
	}
	wg.Wait()
}
