package main

import (
	"bufio"
	"bytes"
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

		go httpHander(conn)
	}
}

func httpHander(conn net.Conn) error {
	defer conn.Close()

	// Wrap the connection in a buffered reader.
	reader := bufio.NewReader(conn)

	req := make([]byte, 0, 4)
	for {
		// starting at the end of the current slice b and extending to its full capacity.
		n, err := reader.Read(req[len(req):cap(req)])
		// resizes the slice b to include the bytes that were just read.
		req = req[:len(req)+n]
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Printf("failed to http handler: %v\n", err)
			return err
		}

		if len(req) == cap(req) {
			// Add more capacity (let append pick how much).
			// and unchanged the original length.
			req = append(req, 0)[:len(req)]
			continue
		}
		break
	}

	headers := make([]string, 0)
	curCont := req
	for {
		buff := bytes.NewBuffer(curCont)
		line, err := buff.ReadString('\n')
		if err != nil {
			log.Printf("failed to http handler: %v\n", err)
			return err
		}

		headers = append(headers, line)
		curCont = curCont[len(line):]
		if line == "\r\n" {
			break
		}
	}

	body := string(curCont)
	log.Println(body)

	log.Println(headers)

	conn.Write([]byte("hello world!"))
	return nil
}
