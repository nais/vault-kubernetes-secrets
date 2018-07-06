package main

import (
	"github.com/nais/vault-kubernetes-secrets/cmd"
	"os"
)

func main() {
	if err := cmd.FetchCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
