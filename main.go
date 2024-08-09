package main

import (
	"fmt"
	"os"

	"github.com/leonardys/pauli-practice/practice"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(practice.New(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
