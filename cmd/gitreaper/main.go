package main

import (
	"fmt"
	"github.com/anasmohammad611/gitreaper/internal/cli"
	"os"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
