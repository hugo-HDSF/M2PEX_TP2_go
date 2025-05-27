package main

import (
	"fmt"
	"os"

	"m2pex_tp2_go/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur lors de l'exécution: %v\n", err)
		os.Exit(1)
	}
}
