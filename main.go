package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/makarellav/ds-pc-kpi/client"
	"github.com/makarellav/ds-pc-kpi/server"
)

func main() {
	t := flag.String("type", "server", "specifies whether to run this program as a tcp server or as a tcp client")
	port := flag.Int("port", 8080, "port for tcp server to listen on")

	flag.Parse()

	switch *t {
	case "server":
		err := server.Listen(*port)

		if err != nil {
			fmt.Printf("failed to start the tcp server: %v", err)

			os.Exit(1)
		}
	case "client":
		err := client.Connect(*port)

		if err != nil {
			fmt.Printf("failed to connect to the tcp server: %v", err)

			os.Exit(1)
		}

	default:
		fmt.Printf("got unknown type: %q, wanted 'client' or 'server'", *t)

		os.Exit(1)
	}
}
