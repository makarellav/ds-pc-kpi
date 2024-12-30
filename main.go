package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/makarellav/ds-pc-kpi/client"
	"github.com/makarellav/ds-pc-kpi/server"
)

func main() {
	mode := flag.String("type", "server", "either 'server' or 'client'")
	port := flag.Int("port", 9000, "port to listen on or connect to")

	flag.Parse()

	switch *mode {
	case "server":
		err := server.Listen(*port)
		if err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
			os.Exit(1)
		}
	case "client":
		err := client.Connect(*port)
		if err != nil {
			fmt.Printf("Failed to connect to server: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown type '%s': expected 'server' or 'client'\n", *mode)
		os.Exit(1)
	}
}
