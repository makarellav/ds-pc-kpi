package client

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// readIntFromStdin reads an integer from stdin.
func readIntFromStdin(prompt string) (int, error) {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return 0, fmt.Errorf("failed to read from stdin")
	}
	text := strings.TrimSpace(scanner.Text())
	value, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("invalid integer: %v", err)
	}
	return value, nil
}

// readMatrixFromStdin reads a rows x cols matrix from user input.
func readMatrixFromStdin(rows, cols int, label string) ([][]int, error) {
	fmt.Printf("Enter elements for %s (%dx%d), row by row:\n", label, rows, cols)
	matrix := make([][]int, rows)
	scanner := bufio.NewScanner(os.Stdin)

	for r := 0; r < rows; r++ {
		rowData := make([]int, cols)
		for c := 0; c < cols; c++ {
			fmt.Printf("Element [%d][%d]: ", r, c)
			if !scanner.Scan() {
				return nil, fmt.Errorf("failed to read matrix row %d", r)
			}
			text := strings.TrimSpace(scanner.Text())
			val, err := strconv.Atoi(text)
			if err != nil {
				return nil, fmt.Errorf("invalid integer in matrix at row %d col %d: %v", r, c, err)
			}
			rowData[c] = val
		}
		matrix[r] = rowData
	}

	return matrix, nil
}

func flattenMatrix(matrix [][]int) []int {
	flat := make([]int, 0, len(matrix)*len(matrix[0]))
	for _, row := range matrix {
		flat = append(flat, row...)
	}
	return flat
}

func sendIntSlice(conn net.Conn, data []int) error {
	buf := make([]byte, 4)
	for _, val := range data {
		binary.BigEndian.PutUint32(buf, uint32(val))
		_, err := conn.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func receiveIntSlice(conn net.Conn, length int) ([]int, error) {
	result := make([]int, length)
	buf := make([]byte, 4)
	for i := 0; i < length; i++ {
		_, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}
		result[i] = int(binary.BigEndian.Uint32(buf))
	}
	return result, nil
}

// Connect prompts the user for matrix dimensions and elements,
// sends them to the server, and receives the multiplication result.
func Connect(port int) error {
	// 1) Read matrix dimensions from user:
	N, err := readIntFromStdin("Enter N (number of rows in matrix A): ")
	if err != nil {
		return fmt.Errorf("error reading N: %v", err)
	}

	M, err := readIntFromStdin("Enter M (number of columns in A / rows in B): ")
	if err != nil {
		return fmt.Errorf("error reading M: %v", err)
	}

	L, err := readIntFromStdin("Enter L (number of columns in matrix B): ")
	if err != nil {
		return fmt.Errorf("error reading L: %v", err)
	}

	// 2) Read matrix A:
	A, err := readMatrixFromStdin(N, M, "Matrix A")
	if err != nil {
		return fmt.Errorf("error reading matrix A: %v", err)
	}

	// 3) Read matrix B:
	B, err := readMatrixFromStdin(M, L, "Matrix B")
	if err != nil {
		return fmt.Errorf("error reading matrix B: %v", err)
	}

	// 4) Dial the server:
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return fmt.Errorf("unable to connect to server: %v", err)
	}
	defer conn.Close()

	buf := make([]byte, 4)

	// 5) Send N, M, L:
	binary.BigEndian.PutUint32(buf, uint32(N))
	conn.Write(buf)
	binary.BigEndian.PutUint32(buf, uint32(M))
	conn.Write(buf)
	binary.BigEndian.PutUint32(buf, uint32(L))
	conn.Write(buf)

	// 6) Send matrix A:
	flatA := flattenMatrix(A)
	if err := sendIntSlice(conn, flatA); err != nil {
		return fmt.Errorf("error sending matrix A: %v", err)
	}

	// 7) Send matrix B:
	flatB := flattenMatrix(B)
	if err := sendIntSlice(conn, flatB); err != nil {
		return fmt.Errorf("error sending matrix B: %v", err)
	}

	// 8) Receive result matrix (N x L):
	resultSize := N * L
	flatResult, err := receiveIntSlice(conn, resultSize)
	if err != nil {
		return fmt.Errorf("error receiving result: %v", err)
	}

	// 9) Print results:
	fmt.Println("Result matrix (flattened) received from server:")
	fmt.Printf("%v\n", flatResult)

	return nil
}
