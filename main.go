package main

import (
	"os"
	"pow-fiat-shamir/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		// Cobra will print the error
		os.Exit(1)
	}
}
