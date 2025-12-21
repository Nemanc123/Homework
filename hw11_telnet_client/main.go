package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var timeout time.Duration

func main() {
	flag.DurationVar(&timeout, "timeout", 60*time.Second, "Connection timeout")
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [-timeout=<time>] <host> <port>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}
	client := NewTelnetClient(fmt.Sprintf("%s:%s", args[0], args[1]), timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		log.Fatal("Connect:", err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = client.Receive()
		err = client.Close()
	}()

	err = client.Send()
	err = client.Close()
	wg.Wait()
}
