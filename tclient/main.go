package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: tclient configfile")
		return
	}
	configs, err := readCofnig(os.Args[1])
	fatalErr(err)
	var wg sync.WaitGroup
	for _, v := range configs {
		wg.Add(1)
		go func(v config) {
			defer wg.Done()
			c := newClient(v)
			go c.run()
			<-c.die
		}(v)
	}
	wg.Wait()
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
