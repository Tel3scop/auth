package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Tel3scop/auth/internal/handlers"
)

const (
	grpcPort = 50051
	protocol = "tcp"
)

func main() {
	lis, err := net.Listen(protocol, fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := handlers.Run()
	log.Printf("server listening at %v", lis.Addr())

	if err = server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
