package main

import (
	"fmt"
	"log"

	"lvbp-project/backend/internal/server"
)

func main() {
	fmt.Println("LBGC Backend starting...")

	srv := server.NewServer()

	if err := srv.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
