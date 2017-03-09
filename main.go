package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
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
