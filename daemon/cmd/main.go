package main

import (
	"fmt"
	"os"
)

func main() {
	svc, err := NewService("config.toml")
	if err != nil {
		fmt.Printf("Error creating service: %v\n", err)
		os.Exit(1)
	}

	svc.Run()

}
