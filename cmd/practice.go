package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	practice "github.com/leonardys/pauli-practice/internal"
	"github.com/spf13/cobra"
)

var practiceCmd = &cobra.Command{
	Use:   "practice",
	Short: "Start a pauli test practice session",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(practice.New(), tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(practiceCmd)
}
