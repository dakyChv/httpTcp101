package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
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
	localAddr := conn.LocalAddr()
	remoteAddr := conn.RemoteAddr()
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("failed to close: %s %s", remoteAddr.Network(), remoteAddr.String())
		}
		log.Printf("%s %s is Closed...\n", remoteAddr.Network(), remoteAddr.String())
	}()
	log.Printf("LocalAddr: %s %s, RemoteAddr(): %s %s\n", localAddr.Network(), localAddr.String(), remoteAddr.Network(), remoteAddr.String())

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

	// Parse the headers and body
	headers := make([]string, 0)
	curCont := req
	for {
		buff := bytes.NewBuffer(curCont)
		line, err := buff.ReadString('\n')
		if err != nil {
			log.Printf("failed to http handler: %v\n", err)
			return err
		}

		curCont = curCont[len(line):]
		if line == "\r\n" {
			break
		}
		headers = append(headers, line)
	}

	headersMap := make(map[string]string, 0)
	for _, v := range headers[1:] {
		kv := strings.Split(v, ":")
		headersMap[kv[0]] = strings.TrimSpace(kv[1])
	}

	// Method, path and version of HTTP.
	// example:
	// Get /mysite/index.html HTTP/1.1\r\n
	// Host: 10.101.101.10\r\n
	// Accept: */*\r\n
	// \r\n
	var (
		method, path, version string
		arrStartLine          []string
	)
	startLine := headers[0]
	arrStartLine = strings.Fields(startLine)
	method = arrStartLine[0]
	path = arrStartLine[1]
	version = arrStartLine[2]

	h, _ := json.Marshal(headersMap)
	log.Println(string(h))
	// print body.
	body := string(curCont)
	log.Println(body)

	// Response.
	// Example:
	//
	// 	HTTP/1.1 200 OK\r\n
	// Content-Length: 55\r\n
	// Content-Type: text/html\r\n
	// Last-Modified: Wed, 12 Aug 1998 15:03:50 GMT\r\n
	// Accept-Ranges: bytes\r\n
	// ETag: “04f97692cbd1:377”\r\n
	// Date: Thu, 19 Jun 2008 19:29:07 GMT\r\n
	// \r\n
	// <55-character response>

	msg := []byte(fmt.Sprintf(`{"request header":%s,"startLine":{"method":%q,"path":%q,"version":%q}}`, h, method, path, version))
	response :=
		"HTTP/1.1 200 OK\r\n" +
			"Content-Type: application/json\r\n" +
			fmt.Sprintf("Content-Length: %d\r\n", len(msg)) +
			"\r\n" +
			string(msg)

	// Reply.
	conn.Write([]byte(response))

	return nil
}
