package server

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"
)

func multiplyMatrices(A, B [][]int) [][]int {
	n := len(A)    // Number of rows in A
	m := len(A[0]) // Number of columns in A (and rows in B)
	l := len(B[0]) // Number of columns in B

	C := make([][]int, n)
	for i := 0; i < n; i++ {
		C[i] = make([]int, l)
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(row int) {
			defer wg.Done()
			for j := 0; j < l; j++ {
				sum := 0
				for k := 0; k < m; k++ {
					sum += A[row][k] * B[k][j]
				}
				C[row][j] = sum
			}
		}(i)
	}
	wg.Wait()
	return C
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 4)

	// Read N
	_, err := conn.Read(buf)
	if err != nil {
		log.Println("Error reading N:", err)
		return
	}
	N := int(binary.BigEndian.Uint32(buf))

	// Read M
	_, err = conn.Read(buf)
	if err != nil {
		log.Println("Error reading M:", err)
		return
	}
	M := int(binary.BigEndian.Uint32(buf))

	// Read L
	_, err = conn.Read(buf)
	if err != nil {
		log.Println("Error reading L:", err)
		return
	}
	L := int(binary.BigEndian.Uint32(buf))

	// Basic dimension check (M must be > 0 if we expect multiplication)
	if M < 1 {
		log.Println("Invalid matrix dimension (M must be > 0)")
		return
	}

	// Read matrix A (N x M)
	A := make([][]int, N)
	for i := 0; i < N; i++ {
		A[i] = make([]int, M)
		for j := 0; j < M; j++ {
			_, err := conn.Read(buf)
			if err != nil {
				log.Println("Error reading matrix A:", err)
				return
			}
			A[i][j] = int(binary.BigEndian.Uint32(buf))
		}
	}

	// Read matrix B (M x L)
	B := make([][]int, M)
	for i := 0; i < M; i++ {
		B[i] = make([]int, L)
		for j := 0; j < L; j++ {
			_, err := conn.Read(buf)
			if err != nil {
				log.Println("Error reading matrix B:", err)
				return
			}
			B[i][j] = int(binary.BigEndian.Uint32(buf))
		}
	}

	// Multiply
	C := multiplyMatrices(A, B)

	// Send result back
	for i := 0; i < N; i++ {
		for j := 0; j < L; j++ {
			binary.BigEndian.PutUint32(buf, uint32(C[i][j]))
			conn.Write(buf)
		}
	}

	fmt.Println("Result sent to client.")
}

func Listen(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen on port %d: %v", port, err)
	}
	defer listener.Close()

	fmt.Printf("Server listening on port %d\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		fmt.Println("New connection...")
		go handleConnection(conn)
	}
}
