package pkg

import (
	"net"
	"tcp_server/pkg/connection"
)

// HandleClient is the entry point for client connection handling
func HandleClient(conn net.Conn) {
	connection.HandleClient(conn)
}
