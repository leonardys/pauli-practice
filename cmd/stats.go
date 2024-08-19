package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	stats "github.com/leonardys/pauli-practice/internal"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "View stats from previous practice sessions",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(stats.NewStatsModel(), tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
