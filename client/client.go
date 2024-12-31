package client

import (
	"encoding/gob"
	"fmt"
	"net"
)

const (
	MessageTypeText   = 1
	MessageTypeBinary = 2
)

type Message struct {
	Type    byte
	Payload []byte
}

func Connect(port int) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Printf("Client connected to port %d\n", port)

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	for i := 1; i <= 100; i++ {
		var msg Message
		if i%2 == 0 {
			msg.Type = MessageTypeText
			msg.Payload = []byte(fmt.Sprintf("Message number %d", i))
		} else {
			msg.Type = MessageTypeBinary
			msg.Payload = []byte{byte(i), byte(i * 2), byte(i * 3)}
		}

		err := encoder.Encode(&msg)
		if err != nil {
			return fmt.Errorf("failed to send message %d: %v", i, err)
		}

		fmt.Printf("Sent message %d of type %d with payload: %v\n", i, msg.Type, msg.Payload)

		var resp Message
		err = decoder.Decode(&resp)
		if err != nil {
			return fmt.Errorf("failed to read response for message %d: %v", i, err)
		}

		fmt.Printf("Received response for message %d: type %d, payload: %v\n", i, resp.Type, resp.Payload)
	}

	return nil
}
