package main

import (
	"fmt"
	"net"
	"tcp_server/pkg"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8181")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

		go pkg.HandleClient(conn)
	}
}
