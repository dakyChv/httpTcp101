package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
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

	for {
		// Waits for connection from listening.
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("failed to next connection to the listener: %w", err)
		}
		defer conn.Close()

		// Wrap the connection in a buffered reader.
		reader := bufio.NewReader(conn)

		b := make([]byte, 0, 4)
		for {
			// starting at the end of the current slice b and extending to its full capacity.
			n, err := reader.Read(b[len(b):cap(b)])
			// resizes the slice b to include the bytes that were just read.
			b = b[:len(b)+n]
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return err
			}

			log.Println("number of bytes read into p: ", n, " and current length of b: ", len(b))
			if len(b) == cap(b) {
				// Add more capacity (let append pick how much).
				// and unchanged the original length.
				b = append(b, 0)[:len(b)]
				log.Println("Extended Cap: ", cap(b))
				continue
			}
			break
		}

		log.Println(string(b))
		conn.Write([]byte("hello world!"))
		return nil
	}
}
