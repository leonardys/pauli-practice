package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pauli-practice",
	Short: "Pauli test preparation",
	Long:  "An application that lets you practice and track your progress on pauli test.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
