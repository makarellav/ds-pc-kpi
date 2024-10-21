package main

import (
	"bytes"
	"fmt"
	"net"
)

func handleConn(conn net.Conn) {
	defer conn.Close()

	data := make([]byte, 1024)

	_, err := conn.Read(data)

	if err != nil {
		fmt.Printf("failed to read data from client: %v", err)

		return
	}

	fmt.Printf("client says: %q\n", bytes.Trim(data, "\x00"))

	hello := []byte("hello, client!")

	fmt.Printf("saying %q to the client\n", hello)

	_, err = conn.Write(hello)

	if err != nil {
		fmt.Printf("failed to write data to client: %v", err)

		return
	}
}

func listen(port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))

	if err != nil {
		return err
	}

	fmt.Printf("server is listening on port %d\n", port)

	defer l.Close()

	for {
		conn, err := l.Accept()

		if err != nil {
			return err
		}

		go handleConn(conn)
	}
}
