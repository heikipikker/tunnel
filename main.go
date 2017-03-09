package main

import (
	"os"
	"fmt"
	"log"
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
		go func() {
			defer wg.Done()
			switch c.Type {
			case "server":
				RunRemoteServer(c)
			default:
				RunLocalServer(c)
			}
		}()
	}
	wg.Wait()
}
