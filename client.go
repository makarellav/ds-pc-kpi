package main

import (
	"bytes"
	"fmt"
	"net"
)

func connect(port int) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))

	if err != nil {
		return err
	}

	fmt.Printf("client is connected to the port %d\n", port)

	defer conn.Close()

	hello := []byte("hello, server!")

	fmt.Printf("saying %q to the server\n", hello)

	_, err = conn.Write(hello)

	if err != nil {
		return err
	}

	data := make([]byte, 1024)

	_, err = conn.Read(data)

	if err != nil {
		return err
	}

	fmt.Printf("server says: %q\n", bytes.Trim(data, "\x00"))

	return nil
}
