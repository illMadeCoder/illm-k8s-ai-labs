package main

import (
	"os"

	"github.com/illmadecoder/labctl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
