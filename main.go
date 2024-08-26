package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	if err := execute(); err != nil {
		log.Fatalf("main(): %v", err)
	}
}

func execute() error {
	// Listen TCP on port 8080.
	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		return fmt.Errorf("failed to listen on the network address: %w", err)
	}
	defer l.Close()
	fmt.Println(`
	/ \__
	(    @\___
	/         O
	/   (_____/
	/_____/   U Server is running at port 8080...`)

	time.Sleep(1 * time.Second)
	return nil
}
