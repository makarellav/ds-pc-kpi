package server

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

func handleConn(conn net.Conn) {
	defer conn.Close()
	decoder := gob.NewDecoder(conn)
	encoder := gob.NewEncoder(conn)

	for i := 0; i < 100; i++ { // Обробка точно 100 повідомлень
		var msg Message
		// Читання повідомлення від клієнта
		err := decoder.Decode(&msg)
		if err != nil {
			fmt.Printf("Failed to read data from client: %v\n", err)
			return
		}

		fmt.Printf("Received message type %d with payload: %v\n", msg.Type, msg.Payload)

		// Обробка повідомлення та формування відповіді
		var response Message
		if msg.Type == MessageTypeText {
			response.Type = MessageTypeText
			response.Payload = []byte("Hello, client!")
		} else if msg.Type == MessageTypeBinary {
			response.Type = MessageTypeBinary
			response.Payload = []byte{0x00, 0x01, 0x02, 0x03}
		} else {
			// Невідомий тип повідомлення
			response.Type = 255
			response.Payload = []byte("Unknown message type")
		}

		// Відправка відповіді клієнту
		err = encoder.Encode(&response)
		if err != nil {
			fmt.Printf("Failed to write data to client: %v\n", err)
			return
		}

		fmt.Printf("Sent response of type %d with payload: %v\n", response.Type, response.Payload)
	}

	fmt.Println("Processed 100 messages, closing connection.")
}

func Listen(port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}

	fmt.Printf("Server is listening on port %d\n", port)
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		go handleConn(conn)
	}
}
