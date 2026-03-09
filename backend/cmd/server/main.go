package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Goldex machine management service")
	fmt.Println("Port gRPC: 50051, HTTP: 8080")
	os.Exit(0)
}
