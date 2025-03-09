package main

import (
	"creditcard/handler"
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println("Incorrect input. Arguments not recognized.")
		os.Exit(1)
	}
	handler.HandlerArgs(args)
}
